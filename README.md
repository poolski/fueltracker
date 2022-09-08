# Fuel Tracker

> This is still a work in progress. I haven't wired up much.

This tool makes calls to the [UK Vehicle Data](https://ukvehicledata.co.uk) Fuel Prices API (for which you will need an API key) to get the 10 nearest fuel stations, and dumps the result to your console.

## Installation

```bash
go get github.com/poolski/fueltracker
go install github.com/poolski/fueltracker
```

## Usage

If you just want to look up the prices of fuel from your console, you'll need to copy the included `config.json.example` to a location on your hard drive and and pass the location to the tool using the `--config` flag.

The config should look like this

```json
{
  "ukvd_api_key": "YOURAPIKEYHERE",
  "google": {
    "credentials_path": "/home/user/.config/fueltracker/service_account.json",
    "spreadsheet_id": "SPREADSHEET_ID_HERE",
    "worksheet_range": "Sheet1!A2"
  }
}
```

Your spreadsheet_id can be found by looking at the address of your Google Sheets spreadsheet, e.g.

```text
https://docs.google.com/spreadsheets/d/SPREADSHEET_ID_HERE/
```

By default, the tool will look for this file at `~/.config/fueltracker/config.json`

If you aren't planning to use Google Sheets, you can omit setting the `google` config.
Otherwise, you will need to follow the instructions [here](https://robocorp.com/docs/development-guide/google-sheets/interacting-with-google-sheets) to generate a service account key for this tool to use.

Make sure that you _share_ the spreadsheet with the Service Account user that you've created, and give it `Editor` access, or it won't be able to write anything to your spreadsheet.

Save the file somewhere on disk and configure the `google.credentials_path` appropriately with the **full path**.

### Usage Examples

```bash
fueltracker --help
fueltracker lookup -p AB123XY
fueltracker write -p AB123XY -f Unleaded -s "STATION NAME"
```

### Dead Man's Snitch

If you want to use this tool on a schedule, e.g. with `systemd` timers (out of scope for this README), you might want to sign up for a free account with [Dead Man's Snitch](https://deadmanssnitch.com), which will tell you if the script fails to run for whatever reason.

You'll need to configure an API key in DMS and create a new Snitch. You will need to add two fields to `~/.config/fueltracker/config.json` to make Fueltracker report back to DMS:

```json
{
  "snitch_api_key": "YOURAPIKEYHERE",
  "snitch_id": "abc123xyz"
}
```
