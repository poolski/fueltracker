# Fuel Tracker

> This is still a work in progress. I haven't wired up much.

This tool makes calls to the [UK Vehicle Data](https://ukvehicledata.co.uk) Fuel Prices API (for which you will need an API key) to get the 10 nearest fuel stations to you which sell _Unleaded_ fuel, and dumps the result to your console.

Yes, it only does Unleaded because only weirdos use anything else.
Also I haven't written the bits that deal with the other fuel types.

I'll convert this to use Cobra/Viper for that CLI goodness later.

I also plan to eventually allow this crappy tool to be run as a Cron job to write rows out to a Google Sheets worksheet so we can all cry in panic at rising fuel prices.

### Usage

```bash
export UKVD_API_KEY=<YOUR_API_KEY_HERE>
go run main.go -postcode AB123XY
```
