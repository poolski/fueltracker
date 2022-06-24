package types

type RawAPIResponse struct {
	BillingAccount struct {
		AccountType      string  `json:"AccountType,omitempty"`
		AccountBalance   float64 `json:"AccountBalance,omitempty"`
		TransactionCost  float64 `json:"TransactionCost,omitempty"`
		ExtraInformation struct {
		} `json:"ExtraInformation,omitempty"`
	} `json:"BillingAccount,omitempty"`
	TechnicalSupport struct {
		ServerID               string        `json:"ServerId,omitempty"`
		RequestID              string        `json:"RequestId,omitempty"`
		QueryDurationMs        int           `json:"QueryDurationMs,omitempty"`
		SupportDate            string        `json:"SupportDate,omitempty"`
		SupportTime            string        `json:"SupportTime,omitempty"`
		SupportCode            string        `json:"SupportCode,omitempty"`
		SupportInformationList []interface{} `json:"SupportInformationList,omitempty"`
	} `json:"TechnicalSupport,omitempty"`
	Request struct {
		RequestGUID     string `json:"RequestGuid,omitempty"`
		PackageID       string `json:"PackageId,omitempty"`
		PackageVersion  int    `json:"PackageVersion,omitempty"`
		ResponseVersion int    `json:"ResponseVersion,omitempty"`
		DataKeys        struct {
			Postcode string `json:"Postcode,omitempty"`
		} `json:"DataKeys,omitempty"`
	} `json:"Request,omitempty"`
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
			HasUnleaded      bool        `json:"HasUnleaded,omitempty"`
			HasSuperUnleaded bool        `json:"HasSuperUnleaded,omitempty"`
			HasDiesel        bool        `json:"HasDiesel,omitempty"`
			HasPremiumDiesel interface{} `json:"HasPremiumDiesel,omitempty"`
			HasLpg           interface{} `json:"HasLpg,omitempty"`
			HasEvCharging    interface{} `json:"HasEvCharging,omitempty"`
		} `json:"Fuel,omitempty"`
		Services struct {
			HasCarWash   interface{} `json:"HasCarWash,omitempty"`
			HasTyrePump  interface{} `json:"HasTyrePump,omitempty"`
			HasWater     interface{} `json:"HasWater,omitempty"`
			HasCashPoint interface{} `json:"HasCashPoint,omitempty"`
			HasCarVacuum interface{} `json:"HasCarVacuum,omitempty"`
		} `json:"Services,omitempty"`
	} `json:"Features,omitempty"`
	FuelPriceCount int `json:"FuelPriceCount,omitempty"`
	FuelPriceList  []struct {
		FuelType            string `json:"FuelType,omitempty"`
		LatestRecordedPrice struct {
			InPence      float64 `json:"InPence,omitempty"`
			InGbp        float64 `json:"InGbp,omitempty"`
			TimeRecorded string  `json:"TimeRecorded,omitempty"`
		} `json:"LatestRecordedPrice,omitempty"`
	} `json:"FuelPriceList,omitempty"`
}

type SpecificFuelPrice struct {
	Station    string
	FuelType   string
	Price      float64
	RecordedAt string
}
