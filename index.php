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
		fwrite(self::$fh, json_encode($object, JSON_PRETTY_PRINT));
		fwrite(self::$fh, "\n");
	}

	public static function console($array) {
		foreach ($array as $line) {
			fwrite(self::$fh, "\t");
			fwrite(self::$fh, $line);
			fwrite(self::$fh, "\n");
		}
	}
}

class Git {
	private static $root = '/srv/http';
	private static $remote = 'origin';
	private static $branch = 'master';

	public static function config($config) {
		if (isset($config->root) && isset($config->remote) && isset($config->branch)) {
			self::$root = $config->root;
			self::$remote = $config->remote;
			self::$branch = $config->branch;
		} else {
			Log::event('Malformed configuration: missing properties in "git" object');
			header('HTTP/1.1 500 Internal Server Error');
			exit(1);
		}
	}

	public static function pull($repository, $branch = NULL, $remote = NULL) {
		$output = array();
		$returnval = 0;
		$branch = (is_null($branch)) ? self::$branch : $branch;
		$remote = (is_null($remote)) ? self::$remote : $remote;
		chdir(self::$root . DIRECTORY_SEPARATOR . $repository);
		exec('git pull ' . escapeshellarg($remote) . ' ' . escapeshellarg($branch), $output, $returnval);
		Log::event('Pulling for ' . $repository);
		Log::console($output);
	}

	public static function search($id, $hooks) {
		$found = 0;
		foreach ($hooks as $hook) {
			if (!isset($hook->id) || !isset($hook->local)) {
				Log::event('Malformed hook object: ' . json_encode($hook));
				header('HTTP/1.1 500 Internal Server Error');
				exit(1);
			}
			if ($hook->id == $id) {
				$branch = (isset($hook->branch)) ? $hook->branch : NULL;
				$remote = (isset($hook->remote)) ? $hook->remote : NULL;
				self::pull($hook->local, $branch, $remote);
				$found++;
			}
		}
		return $found;
	}	
}

Log::open($log);

$input = json_decode(file_get_contents('php://input'));
$conf = json_decode(file_get_contents($config));

if (isset($conf->git)) {
	Git::config($conf->git);
} else {
	Log::event('Malformed configuration: missing "git" object');
	header('HTTP/1.1 500 Internal Server Error');
	exit(1);
}

if (array_key_exists('HTTP_X_GITHUB_EVENT', $_SERVER) && isset($input->repository->id)) { // GitHub
	if (!isset($conf->github)) {
		Log::event('Missing configuration to handle GitHub requests');
		header('HTTP/1.1 500 Internal Server Error');
		exit(1);
	}
	if ($_SERVER['HTTP_X_GITHUB_EVENT'] == 'push') {
		if(! Git::search($input->repository->id, $conf->github)) {
			Log::event('Received GitHub pull webhook for an unhandled project');
			Log::struct($input);
		}
	} else {
		Log::event('Received GitHub webhook of unhandled kind ('. $_SERVER['HTTP_X_GITHUB_EVENT'] .')');
		Log::struct($input);
	}
} elseif (isset($input->project_id)) { // GitLab
	if (!isset($conf->gitlab)) {
		Log::event('Missing configuration to handle GitLab requests');
		header('HTTP/1.1 500 Internal Server Error');
		exit(1);
	}
	if (isset($input->object_kind) && $input->object_kind == 'push') {
		if (! Git::search($input->project_id, $conf->gitlab)) {
			Log::event('Received GitLab pull webhook for an unhandled project');
			Log::struct($input);
		}
	} else {
		Log::event('Received GitLab webhook of unhandled kind');
		Log::struct($input);
	}
} else { // Something else
	Log::event('Received unhandled request');
	Log::struct($input);
}

Log::close();
?>
