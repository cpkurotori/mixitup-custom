# usercount

`usercount` command is a utility for a counter per user (based on MixItUp GUID).

## Usage
```
Usage:
  mixitup-custom user-count [flags]

Flags:
      --add int               the amount to add (default 1)
      --counter-file string   file to store counters in (required); will create this file if it does not exist
  -h, --help                  help for user-count
      --user-id string        user id to add counter to (required)
      --user-name string      user name to attach to id

Global Flags:
      --log-file -         log file, - for StdErr (default "-")
      --log-level string   log level [DEBUG|INFO|WARN|ERROR] (default "WARN")
```

## Setup
1. Download the release appropriate for your platform
1. Create a folder where the executable and relevant data will be stored (example: `C:\Users\johndoe\Documents\mixitup-custom`)
1. Move the downloaded executable to this folder
1. Download the [checkin_actions.miucommand](/checkin_actions.miucommand)
1. Open the file and replace the following:
    - `${REPLACE WITH PATH TO EXECUTABLE}` with the full path to your executable (example: `C:\Users\johndoe\Documents\mixitup-custom\mixitup-custom-windows-x86.exe` )
    - `${REPLACE WITH PATH TO EXECUTABLE FOLDER}` with the full path to your folder you created (example: `C:\Users\johndoe\Documents\mixitup-custom` )
1. Open MixItUp and navigate to `Action Groups`
1. Click `New Action Group`
1. At the bottom, click the right most icon (`Import Actions From File`)
1. Select the `checkin_actions.miucommand` file you downloaded and edited
1. Click the `Test Command` (play) button a couple times and validate that the counter increases in your Twitch Chat
1. Save the command and use the Action Group as you wish