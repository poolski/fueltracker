package types

type RawAPIResponse struct {
	Response FuelDataResponse `json:"Response,omitempty"`
}

type FuelDataResponse struct {
	StatusCode        string `json:"StatusCode,omitempty"`
	StatusMessage     string `json:"StatusMessage,omitempty"`
	StatusInformation struct {
		Lookup struct {
			StatusCode    string `json:"StatusCode,omitempty"`
			StatusMessage string `json:"StatusMessage,omitempty"`
		} `json:"Lookup,omitempty"`
	} `json:"StatusInformation,omitempty"`
	DataItems struct {
		FuelStationDetails struct {
			FuelStationCount int           `json:"FuelStationCount,omitempty"`
			SearchRadiusUsed int           `json:"SearchRadiusUsed,omitempty"`
			FuelStationList  []FuelStation `json:"FuelStationList,omitempty"`
		} `json:"FuelStationDetails,omitempty"`
	} `json:"DataItems,omitempty"`
}

type FuelStation struct {
	DistanceFromSearchPostcode float64 `json:"DistanceFromSearchPostcode,omitempty"`
	Brand                      string  `json:"Brand,omitempty"`
	Name                       string  `json:"Name,omitempty"`
	Street                     string  `json:"Street,omitempty"`
	Suburb                     string  `json:"Suburb,omitempty"`
	Town                       string  `json:"Town,omitempty"`
	County                     string  `json:"County,omitempty"`
	Postcode                   string  `json:"Postcode,omitempty"`
	Latitude                   float64 `json:"Latitude,omitempty"`
	Longitude                  float64 `json:"Longitude,omitempty"`
	Features                   struct {
		Fuel struct {
			HasUnleaded      bool `json:"HasUnleaded,omitempty"`
			HasSuperUnleaded bool `json:"HasSuperUnleaded,omitempty"`
			HasDiesel        bool `json:"HasDiesel,omitempty"`
			HasPremiumDiesel bool `json:"HasPremiumDiesel,omitempty"`
			HasLpg           bool `json:"HasLpg,omitempty"`
			HasEvCharging    bool `json:"HasEvCharging,omitempty"`
		} `json:"Fuel,omitempty"`
		Services struct {
			HasCarWash   bool `json:"HasCarWash,omitempty"`
			HasTyrePump  bool `json:"HasTyrePump,omitempty"`
			HasWater     bool `json:"HasWater,omitempty"`
			HasCashPoint bool `json:"HasCashPoint,omitempty"`
			HasCarVacuum bool `json:"HasCarVacuum,omitempty"`
		} `json:"Services,omitempty"`
	} `json:"Features,omitempty"`
	FuelPriceCount int         `json:"FuelPriceCount,omitempty"`
	FuelPriceList  []FuelPrice `json:"FuelPriceList,omitempty"`
}

type FuelPrice struct {
	FuelType            string `json:"FuelType,omitempty"`
	LatestRecordedPrice struct {
		InPence      float64 `json:"InPence,omitempty"`
		InGbp        float64 `json:"InGbp,omitempty"`
		TimeRecorded string  `json:"TimeRecorded,omitempty"`
	} `json:"LatestRecordedPrice,omitempty"`
}
type SpecificFuelPrice struct {
	Station    string
	FuelType   string
	Price      float64
	RecordedAt string
}
