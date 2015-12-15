<?php
/*****************************************************************
 * Dirty Continuous Integration
 * Receives GitHub and GitLab webhooks for selected projects
 * and git pull their content into specified directory.
 * Specifically designed for web projects on a web server...
 * 
 * @author Rémi Dubois ( GitHub @reveliant / GitLab @packman )
 * @license MIT
**/

$config = "config.json";
$log = "hooks.log";

/****************************************************************/

class Log {
	private static $fh = NULL;

	public static function open($path) {
		self::$fh = fopen($path, 'a');
	}

	public static function close() {
		fclose(self::$fh);
	}

	public static function event($message) {
		fwrite(self::$fh, strftime('%F %T'));
		fwrite(self::$fh, "\t");
		fwrite(self::$fh, $message);
		fwrite(self::$fh, "\n");
	}

	public static function struct($object) {
		if ($object !== null) {
			fwrite(self::$fh, json_encode($object, JSON_PRETTY_PRINT));
			fwrite(self::$fh, "\n");
		}
	}

	public static function console($array) {
		foreach ($array as $line) {
			fwrite(self::$fh, "\t");
			fwrite(self::$fh, $line);
			fwrite(self::$fh, "\n");
		}
	}

	public static function client($message) {
		print($message);
	}
}

class Git {
	private static $root = '/srv/http';
	private static $remoteBranch = 'origin';
	private static $localBranch = 'master';
	private static $repositories = Array();

	public static function config($config) {
		if (isset($config->root) && isset($config->branches) && isset($config->branches->remote) && isset($config->branches->local)) {
			self::$root = $config->root;
			self::$remoteBranch = $config->branches->remote;
			self::$localBranch = $config->branches->local;
		} else {
			header('HTTP/1.1 500 Internal Server Error');
			Log::event('Malformed configuration: missing properties in "git" object');
			Log::client('Malformed configuration on server');
			exit(1);
		}
	}

	public static function repositories($repositories) {
		self::$repositories = $repositories;
	}

	public static function pull($repository, $localBranch = NULL, $remoteBranch = NULL) {
		$output = array();
		$returnval = 0;
		$localBranch = (is_null($localBranch)) ? self::$localBranch : $localBranch;
		$remoteBranch = (is_null($remoteBranch)) ? self::$remoteBranch : $remoteBranch;
		chdir(self::$root . DIRECTORY_SEPARATOR . $repository);
		exec('git pull ' . escapeshellarg($remoteBranch) . ' ' . escapeshellarg($localBranch), $output, $returnval);
		Log::event('Pulling for ' . $repository);
		Log::console($output);
	}

	public static function search($url) {
		$found = 0;
		foreach (self::$repositories as $hook) {
			if (!isset($hook->remote) || !isset($hook->remote->url) || !isset($hook->local) || !isset($hook->local->url)) {
				Log::event('Malformed hook object: ' . json_encode($hook));
				break;
			}
			if ($hook->remote->url == $url) {
				$localBranch = (isset($hook->local->branch)) ? $hook->branch->local : NULL;
				$remoteBranch = (isset($hook->remote->branch)) ? $hook->branch->remote : NULL;
				self::pull($hook->local->url, $localBranch, $remoteBranch);
				$found++;
			}
		}
		if ($found) {
	  		header('HTTP/1.1 204 No Content'); 
		} else {
	  		header('HTTP/1.1 202 Accepted'); 
			Log::event('Received a pull webhook for an unhandled project');
			Log::client('This repository may not be synchronized');
		}
		return $found;
	}	
}

Log::open($log);

$input = json_decode(file_get_contents('php://input'));
$conf = json_decode(file_get_contents($config));

if (isset($conf->git) && isset($conf->repositories)) {
	Git::config($conf->git);
	Git::repositories($conf->repositories);
} else {
	header('HTTP/1.1 500 Internal Server Error');
	Log::event('Malformed configuration');
	exit(1);
}

if (array_key_exists('HTTP_X_GITHUB_EVENT', $_SERVER) && isset($input->repository->ssh_url)) { // GitHub
	if ($_SERVER['HTTP_X_GITHUB_EVENT'] == 'push') {
		if(! Git::search($input->repository->ssh_url)) {
			Log::struct($input);
		}
	} else {
		Log::event('Received a GitHub webhook of unhandled kind ('. $_SERVER['HTTP_X_GITHUB_EVENT'] .')');
		Log::struct($input);
	}
} elseif (array_key_exists('HTTP_X_GITLAB_EVENT', $_SERVER) && isset($input->repository->git_ssh_url)) { // GitLab
	if ($_SERVER['HTTP_X_GITLAB_EVENT'] == 'Push Hook') {
		if (! Git::search($input->repository->git_ssh_url)) {
			Log::struct($input);
		}
	} else {
		Log::event('Received a GitLab webhook of unhandled kind ('. $_SERVER['HTTP_X_GITLAB_EVENT'] . ')');
		Log::struct($input);
	}
} else { // Something else
	header('HTTP/1.1 501 Not Implemented'); 
	Log::client('Unhandled request');
	Log::event('Received unhandled request');
	Log::struct($input);
}

Log::close();
?>
