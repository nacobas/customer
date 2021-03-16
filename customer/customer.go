package customer

import (
	"regexp"
	"strings"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/pkg/errors"
)

var (
	ErrInvalidValue = errors.New("Invalid value")
	ErrUnkownType   = errors.New("Unknown type")
)

func New(id uint64, info CustomerInfo) (Customer, error) {
	switch i := info.(type) {
	case PersonInfo:
		ret, err := newPrivate(id, i)
		if err != nil {
			return nil, errors.Wrap(err, "customer.New")
		}
		return ret, nil

	case OrganizationInfo:
		ret, err := newOrganization(id, i)
		if err != nil {
			return nil, errors.Wrap(err, "customer.New")
		}
		return ret, nil

	default:
		return nil, errors.Wrap(ErrUnkownType, "customer.New")
	}
}

type Customer interface {
	ID() uint64
	Info() CustomerInfo
	Addresses() []*Address
	ContactInfos() []ContactInfo
	TaxInfos() []*TaxInfo
}

func newPrivate(id uint64, pi PersonInfo) (*private, error) {
	return &private{id: id, PersonInfo: pi}, nil
}

type private struct {
	id uint64
	PersonInfo
	addrs    []*Address
	contacts []ContactInfo
	taxs     []*TaxInfo
}

func (p *private) ID() uint64 {
	return p.id
}

func (p *private) Info() CustomerInfo {
	return p.PersonInfo
}

func (p *private) Addresses() []*Address {
	return p.addrs
}

func (p *private) ContactInfos() []ContactInfo {
	return p.contacts
}

func (p *private) TaxInfos() []*TaxInfo {
	return p.taxs
}

func newOrganization(id uint64, oi OrganizationInfo) (*organization, error) {
	return &organization{id: id, OrganizationInfo: oi}, nil
}

type organization struct {
	id uint64
	OrganizationInfo
	addrs    []*Address
	contacts []ContactInfo
	taxs     []*TaxInfo
}

func (o *organization) ID() uint64 {
	return o.id
}

func (o *organization) Info() CustomerInfo {
	return o.OrganizationInfo
}

func (o *organization) Addresses() []*Address {
	return o.addrs
}

func (o *organization) ContactInfos() []ContactInfo {
	return o.contacts
}

func (o *organization) TaxInfos() []*TaxInfo {
	return o.taxs
}

type CustomerInfo interface {
	IsCustomerInfo()
}

type PersonInfo struct {
	GivenName   string
	FamilyName  string
	SSN         string
	DateOfBirth date.Date
	Citizenship CountryCode
}

func (p PersonInfo) IsCustomerInfo()

type OrganizationInfo struct {
	Name                string
	Form                string
	LeagalID            string
	RegistrationDate    date.Date
	RegistrationCountry CountryCode
}

func (o OrganizationInfo) IsCustomerInfo()

type Address interface {
	Street() string
	City() string
	PostalCode() string
	Country() string
}

type address struct {
	street   string
	city     string
	postCode string
	cc       CountryCode
}

func (a *address) Street() string {
	return a.street
}
func (a *address) City() string {
	return a.city
}

func (a *address) PostalCode() string {
	return a.postCode
}

func (a *address) Country() string {
	return a.cc.Country()
}

func NewContactInfo(ct ContactType, value string) (ContactInfo, error) {
	switch ct {
	case ContactTypeMobile:
		return Mobile{value}, nil
	case ContactTypePhone:
		return Phone{value}, nil
	case ContactTypeEmail:
		return Email{value}, nil
	default:
		return nil, errors.Wrap(ErrUnkownType, "customer.NewContactInfo")
	}
}

type ContactInfo interface {
	Value() string
	Type() ContactType
}

type ContactType int

const (
	ContactTypeOther ContactType = iota
	ContactTypeMobile
	ContactTypePhone
	ContactTypeEmail
)

var phoneRegex = regexp.MustCompile(`^\+?([0-9]+)$`)

func NewMobile(number string) (Mobile, error) {
	s := strings.Trim(number, " ")

	ok := phoneRegex.MatchString(s)
	if !ok {
		return Mobile{}, errors.Wrap(ErrInvalidValue, "customer.NewMobile")
	}

	return Mobile{Number: s}, nil
}

type Mobile struct {
	Number string
}

func (m Mobile) Value() string {
	return m.Number
}

func (m Mobile) Type() ContactType {
	return ContactTypeMobile
}

func (m Mobile) Validate() error {
	ok := phoneRegex.MatchString(m.Number)
	if !ok {
		return errors.Wrap(ErrInvalidValue, "customer.Mobile.Validate")
	}
	return nil

}

func NewPhone(number string) (Phone, error) {
	ok := phoneRegex.MatchString(number)
	if !ok {
		return Phone{}, errors.Wrap(ErrInvalidValue, "customer.NewPhone")
	}
	return Phone{Number: number}, nil
}

type Phone struct {
	Number string
}

func (p Phone) Value() string {
	return p.Number
}

func (p Phone) Type() ContactType {
	return ContactTypePhone
}

func (p Phone) Validate() error {
	ok := phoneRegex.MatchString(p.Number)
	if !ok {
		return errors.Wrap(ErrInvalidValue, "customer.Phone.Validate")
	}
	return nil

}

var emailRegexp = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")

func NewEmail(address string) (Email, error) {
	ok := emailRegexp.MatchString(address)
	if !ok {
		return Email{}, errors.Wrap(ErrInvalidValue, "customer.NewEmail")
	}
	return Email{Address: address}, nil
}

type Email struct {
	Address string
}

func (e Email) Value() string {
	return e.Address
}

func (e Email) Type() ContactType {
	return ContactTypePhone
}

func (e Email) Validate() error {
	ok := emailRegexp.MatchString(e.Address)
	if !ok {
		return errors.Wrap(ErrInvalidValue, "customer.Email.Validate")
	}
	return nil
}

type TaxInfo struct {
	TaxCountry CountryCode
	TaxId      string
}

type CountryCode int32

func (cc CountryCode) String() string {
	return cc_code[cc]
}

func (cc CountryCode) Country() string {
	return cc_country[cc]
}

const (
	ABW CountryCode = iota + 1 // Aruba
	AFG                        // Afghanistan
	AGO                        // Angola
	AIA                        // Anguilla
	ALA                        // Åland Islands
	ALB                        // Albania
	AND                        // Andorra
	ANT                        // Netherlands Antilles
	ARE                        // United Arab Emirates
	ARG                        // Argentina
	ARM                        // Armenia
	ASM                        // American Samoa
	ATA                        // Antarctica
	ATF                        // French Southern Territories
	ATG                        // Antigua and Barbuda
	AUS                        // Australia
	AUT                        // Austria
	AZE                        // Azerbaijan
	BDI                        // Burundi
	BEL                        // Belgium
	BEN                        // Benin
	BFA                        // Burkina Faso
	BGD                        // Bangladesh
	BGR                        // Bulgaria
	BHR                        // Bahrain
	BHS                        // Bahamas
	BIH                        // Bosnia and Herzegovina
	BLM                        // Saint Barthélemy
	BLR                        // Belarus
	BLZ                        // Belize
	BMU                        // Bermuda
	BOL                        // Bolivia
	BRA                        // Brazil
	BRB                        // Barbados
	BRN                        // Brunei Darussalam
	BTN                        // Bhutan
	BVT                        // Bouvet Island
	BWA                        // Botswana
	CAF                        // Central African Republic
	CAN                        // Canada
	CCK                        // Cocos (Keeling) Islands
	CHE                        // Switzerland
	CHL                        // Chile
	CHN                        // China
	CIV                        // Côte d'Ivoire
	CMR                        // Cameroon
	COD                        // Congo, the Democratic Republic of the
	COG                        // Congo
	COK                        // Cook Islands
	COL                        // Colombia
	COM                        // Comoros
	CPV                        // Cape Verde
	CRI                        // Costa Rica
	CUB                        // Cuba
	CXR                        // Christmas Island
	CYM                        // Cayman Islands
	CYP                        // Cyprus
	CZE                        // Czechia
	DEU                        // Germany
	DJI                        // Djibouti
	DMA                        // Dominica
	DNK                        // Denmark
	DOM                        // Dominican Republic
	DZA                        // Algeria
	ECU                        // Ecuador
	EGY                        // Egypt
	ERI                        // Eritrea
	ESH                        // Western Sahara
	ESP                        // Spain
	EST                        // Estonia
	ETH                        // Ethiopia
	FIN                        // Finland
	FJI                        // Fiji
	FLK                        // Falkland Islands (Malvinas)
	FRA                        // France
	FRO                        // Faroe Islands
	FSM                        // Micronesia, Federated States of
	GAB                        // Gabon
	GBR                        // United Kingdom
	GEO                        // Georgia
	GGY                        // Guernsey
	GHA                        // Ghana
	GIN                        // Guinea
	GIB                        // Gibraltar
	GLP                        // Guadeloupe
	GMB                        // Gambia
	GNB                        // Guinea-Bissau
	GNQ                        // Equatorial Guinea
	GRC                        // Greece
	GRD                        // Grenada
	GRL                        // Greenland
	GTM                        // Guatemala
	GUF                        // French Guiana
	GUM                        // Guam
	GUY                        // Guyana
	HKG                        // Hong Kong
	HMD                        // Heard Island and McDonald Islands
	HND                        // Honduras
	HRV                        // Croatia
	HTI                        // Haiti
	HUN                        // Hungary
	IDN                        // Indonesia
	IMN                        // Isle of Man
	IND                        // India
	IOT                        // British Indian Ocean Territory
	IRL                        // Ireland
	IRN                        // Iran, Islamic Republic of
	IRQ                        // Iraq
	ISL                        // Iceland
	ISR                        // Israel
	ITA                        // Italy
	JAM                        // Jamaica
	JEY                        // Jersey
	JOR                        // Jordan
	JPN                        // Japan
	KAZ                        // Kazakhstan
	KEN                        // Kenya
	KGZ                        // Kyrgyzstan
	KHM                        // Cambodia
	KIR                        // Kiribati
	KNA                        // Saint Kitts and Nevis
	KOR                        // Korea, Republic of
	KWT                        // Kuwait
	LAO                        // Lao People's Democratic Republic
	LBN                        // Lebanon
	LBR                        // Liberia
	LBY                        // Libyan Arab Jamahiriya
	LCA                        // Saint Lucia
	LIE                        // Liechtenstein
	LKA                        // Sri Lanka
	LSO                        // Lesotho
	LTU                        // Lithuania
	LUX                        // Luxembourg
	LVA                        // Latvia
	MAC                        // Macao
	MAF                        // Saint Martin (French part)
	MAR                        // Morocco
	MCO                        // Monaco
	MDA                        // Moldova, Republic of
	MDG                        // Madagascar
	MDV                        // Maldives
	MEX                        // Mexico
	MHL                        // Marshall Islands
	MKD                        // Macedonia, the former Yugoslav Republic of
	MLI                        // Mali
	MLT                        // Malta
	MMR                        // Myanmar
	MNE                        // Montenegro
	MNG                        // Mongolia
	MNP                        // Northern Mariana Islands
	MOZ                        // Mozambique
	MRT                        // Mauritania
	MSR                        // Montserrat
	MTQ                        // Martinique
	MUS                        // Mauritius
	MWI                        // Malawi
	MYS                        // Malaysia
	MYT                        // Mayotte
	NAM                        // Namibia
	NCL                        // New Caledonia
	NER                        // Niger
	NFK                        // Norfolk Island
	NGA                        // Nigeria
	NIC                        // Nicaragua
	NOR                        // Norway
	NIU                        // Niue
	NLD                        // Netherlands
	NPL                        // Nepal
	NRU                        // Nauru
	NZL                        // New Zealand
	OMN                        // Oman
	PAK                        // Pakistan
	PAN                        // Panama
	PCN                        // Pitcairn
	PER                        // Peru
	PHL                        // Philippines
	PLW                        // Palau
	PNG                        // Papua New Guinea
	POL                        // Poland
	PRI                        // Puerto Rico
	PRK                        // Korea, Democratic People's Republic of
	PRT                        // Portugal
	PRY                        // Paraguay
	PSE                        // Palestinian Territory, Occupied
	PYF                        // French Polynesia
	QAT                        // Qatar
	REU                        // Réunion
	ROU                        // Romania
	RUS                        // Russian Federation
	RWA                        // Rwanda
	SAU                        // Saudi Arabia
	SDN                        // Sudan
	SEN                        // Senegal
	SGP                        // Singapore
	SGS                        // South Georgia and the South Sandwich Islands
	SHN                        // Saint Helena
	SJM                        // Svalbard and Jan Mayen
	SLB                        // Solomon Islands
	SLE                        // Sierra Leone
	SLV                        // El Salvador
	SMR                        // San Marino
	SOM                        // Somalia
	SPM                        // Saint Pierre and Miquelon
	SRB                        // Serbia
	STP                        // São Tomé and Príncipe
	SUR                        // Suriname
	SVK                        // Slovakia
	SVN                        // Slovenia
	SWE                        // Sweden
	SWZ                        // Swaziland
	SYC                        // Seychelles
	SYR                        // Syrian Arab Republic
	TCA                        // Turks and Caicos Islands
	TCD                        // Chad
	TGO                        // Togo
	THA                        // Thailand
	TJK                        // Tajikistan
	TKL                        // Tokelau
	TKM                        // Turkmenistan
	TLS                        // Timor-Leste
	TON                        // Tonga
	TTO                        // Trinidad and Tobago
	TUN                        // Tunisia
	TUR                        // Turkey
	TUV                        // Tuvalu
	TWN                        // Taiwan, Province of China
	TZA                        // Tanzania, United Republic of
	UGA                        // Uganda
	UKR                        // Ukraine
	UMI                        // United States Minor Outlying Islands
	URY                        // Uruguay
	USA                        // United States
	UZB                        // Uzbekistan
	VAT                        // Holy See (Vatican City State)
	VCT                        // Saint Vincent and the Grenadines
	VEN                        // Venezuela
	VGB                        // Virgin Islands, British
	VIR                        // Virgin Islands, U.S.
	VNM                        // Viet Nam
	VUT                        // Vanuatu
	WLF                        // Wallis and Futuna
	WSM                        // Samoa
	YEM                        // Yemen
	ZAF                        // South Africa
	ZMB                        // Zambia
	ZWE                        // Zimbabwe
)

// Enum value maps for CountryCode.
var (
	cc_code = map[CountryCode]string{
		1:   "ABW",
		2:   "AFG",
		3:   "AGO",
		4:   "AIA",
		5:   "ALA",
		6:   "ALB",
		7:   "AND",
		8:   "ANT",
		9:   "ARE",
		10:  "ARG",
		11:  "ARM",
		12:  "ASM",
		13:  "ATA",
		14:  "ATF",
		15:  "ATG",
		16:  "AUS",
		17:  "AUT",
		18:  "AZE",
		19:  "BDI",
		20:  "BEL",
		21:  "BEN",
		22:  "BFA",
		23:  "BGD",
		24:  "BGR",
		25:  "BHR",
		26:  "BHS",
		27:  "BIH",
		28:  "BLM",
		29:  "BLR",
		30:  "BLZ",
		31:  "BMU",
		32:  "BOL",
		33:  "BRA",
		34:  "BRB",
		35:  "BRN",
		36:  "BTN",
		37:  "BVT",
		38:  "BWA",
		39:  "CAF",
		40:  "CAN",
		41:  "CCK",
		42:  "CHE",
		43:  "CHL",
		44:  "CHN",
		45:  "CIV",
		46:  "CMR",
		47:  "COD",
		48:  "COG",
		49:  "COK",
		50:  "COL",
		51:  "COM",
		52:  "CPV",
		53:  "CRI",
		54:  "CUB",
		55:  "CXR",
		56:  "CYM",
		57:  "CYP",
		58:  "CZE",
		59:  "DEU",
		60:  "DJI",
		61:  "DMA",
		62:  "DNK",
		63:  "DOM",
		64:  "DZA",
		65:  "ECU",
		66:  "EGY",
		67:  "ERI",
		68:  "ESH",
		69:  "ESP",
		70:  "EST",
		71:  "ETH",
		72:  "FIN",
		73:  "FJI",
		74:  "FLK",
		75:  "FRA",
		76:  "FRO",
		77:  "FSM",
		78:  "GAB",
		79:  "GBR",
		80:  "GEO",
		81:  "GGY",
		82:  "GHA",
		83:  "GIN",
		84:  "GIB",
		85:  "GLP",
		86:  "GMB",
		87:  "GNB",
		88:  "GNQ",
		89:  "GRC",
		90:  "GRD",
		91:  "GRL",
		92:  "GTM",
		93:  "GUF",
		94:  "GUM",
		95:  "GUY",
		96:  "HKG",
		97:  "HMD",
		98:  "HND",
		99:  "HRV",
		100: "HTI",
		101: "HUN",
		102: "IDN",
		103: "IMN",
		104: "IND",
		105: "IOT",
		106: "IRL",
		107: "IRN",
		108: "IRQ",
		109: "ISL",
		110: "ISR",
		111: "ITA",
		112: "JAM",
		113: "JEY",
		114: "JOR",
		115: "JPN",
		116: "KAZ",
		117: "KEN",
		118: "KGZ",
		119: "KHM",
		120: "KIR",
		121: "KNA",
		122: "KOR",
		123: "KWT",
		124: "LAO",
		125: "LBN",
		126: "LBR",
		127: "LBY",
		128: "LCA",
		129: "LIE",
		130: "LKA",
		131: "LSO",
		132: "LTU",
		133: "LUX",
		134: "LVA",
		135: "MAC",
		136: "MAF",
		137: "MAR",
		138: "MCO",
		139: "MDA",
		140: "MDG",
		141: "MDV",
		142: "MEX",
		143: "MHL",
		144: "MKD",
		145: "MLI",
		146: "MLT",
		147: "MMR",
		148: "MNE",
		149: "MNG",
		150: "MNP",
		151: "MOZ",
		152: "MRT",
		153: "MSR",
		154: "MTQ",
		155: "MUS",
		156: "MWI",
		157: "MYS",
		158: "MYT",
		159: "NAM",
		160: "NCL",
		161: "NER",
		162: "NFK",
		163: "NGA",
		164: "NIC",
		165: "NOR",
		166: "NIU",
		167: "NLD",
		168: "NPL",
		169: "NRU",
		170: "NZL",
		171: "OMN",
		172: "PAK",
		173: "PAN",
		174: "PCN",
		175: "PER",
		176: "PHL",
		177: "PLW",
		178: "PNG",
		179: "POL",
		180: "PRI",
		181: "PRK",
		182: "PRT",
		183: "PRY",
		184: "PSE",
		185: "PYF",
		186: "QAT",
		187: "REU",
		188: "ROU",
		189: "RUS",
		190: "RWA",
		191: "SAU",
		192: "SDN",
		193: "SEN",
		194: "SGP",
		195: "SGS",
		196: "SHN",
		197: "SJM",
		198: "SLB",
		199: "SLE",
		200: "SLV",
		201: "SMR",
		202: "SOM",
		203: "SPM",
		204: "SRB",
		205: "STP",
		206: "SUR",
		207: "SVK",
		208: "SVN",
		209: "SWE",
		210: "SWZ",
		211: "SYC",
		212: "SYR",
		213: "TCA",
		214: "TCD",
		215: "TGO",
		216: "THA",
		217: "TJK",
		218: "TKL",
		219: "TKM",
		220: "TLS",
		221: "TON",
		222: "TTO",
		223: "TUN",
		224: "TUR",
		225: "TUV",
		226: "TWN",
		227: "TZA",
		228: "UGA",
		229: "UKR",
		230: "UMI",
		231: "URY",
		232: "USA",
		233: "UZB",
		234: "VAT",
		235: "VCT",
		236: "VEN",
		237: "VGB",
		238: "VIR",
		239: "VNM",
		240: "VUT",
		241: "WLF",
		242: "WSM",
		243: "YEM",
		244: "ZAF",
		245: "ZMB",
		246: "ZWE",
	}
	cc_value = map[string]CountryCode{
		"ABW": 1,
		"AFG": 2,
		"AGO": 3,
		"AIA": 4,
		"ALA": 5,
		"ALB": 6,
		"AND": 7,
		"ANT": 8,
		"ARE": 9,
		"ARG": 10,
		"ARM": 11,
		"ASM": 12,
		"ATA": 13,
		"ATF": 14,
		"ATG": 15,
		"AUS": 16,
		"AUT": 17,
		"AZE": 18,
		"BDI": 19,
		"BEL": 20,
		"BEN": 21,
		"BFA": 22,
		"BGD": 23,
		"BGR": 24,
		"BHR": 25,
		"BHS": 26,
		"BIH": 27,
		"BLM": 28,
		"BLR": 29,
		"BLZ": 30,
		"BMU": 31,
		"BOL": 32,
		"BRA": 33,
		"BRB": 34,
		"BRN": 35,
		"BTN": 36,
		"BVT": 37,
		"BWA": 38,
		"CAF": 39,
		"CAN": 40,
		"CCK": 41,
		"CHE": 42,
		"CHL": 43,
		"CHN": 44,
		"CIV": 45,
		"CMR": 46,
		"COD": 47,
		"COG": 48,
		"COK": 49,
		"COL": 50,
		"COM": 51,
		"CPV": 52,
		"CRI": 53,
		"CUB": 54,
		"CXR": 55,
		"CYM": 56,
		"CYP": 57,
		"CZE": 58,
		"DEU": 59,
		"DJI": 60,
		"DMA": 61,
		"DNK": 62,
		"DOM": 63,
		"DZA": 64,
		"ECU": 65,
		"EGY": 66,
		"ERI": 67,
		"ESH": 68,
		"ESP": 69,
		"EST": 70,
		"ETH": 71,
		"FIN": 72,
		"FJI": 73,
		"FLK": 74,
		"FRA": 75,
		"FRO": 76,
		"FSM": 77,
		"GAB": 78,
		"GBR": 79,
		"GEO": 80,
		"GGY": 81,
		"GHA": 82,
		"GIN": 83,
		"GIB": 84,
		"GLP": 85,
		"GMB": 86,
		"GNB": 87,
		"GNQ": 88,
		"GRC": 89,
		"GRD": 90,
		"GRL": 91,
		"GTM": 92,
		"GUF": 93,
		"GUM": 94,
		"GUY": 95,
		"HKG": 96,
		"HMD": 97,
		"HND": 98,
		"HRV": 99,
		"HTI": 100,
		"HUN": 101,
		"IDN": 102,
		"IMN": 103,
		"IND": 104,
		"IOT": 105,
		"IRL": 106,
		"IRN": 107,
		"IRQ": 108,
		"ISL": 109,
		"ISR": 110,
		"ITA": 111,
		"JAM": 112,
		"JEY": 113,
		"JOR": 114,
		"JPN": 115,
		"KAZ": 116,
		"KEN": 117,
		"KGZ": 118,
		"KHM": 119,
		"KIR": 120,
		"KNA": 121,
		"KOR": 122,
		"KWT": 123,
		"LAO": 124,
		"LBN": 125,
		"LBR": 126,
		"LBY": 127,
		"LCA": 128,
		"LIE": 129,
		"LKA": 130,
		"LSO": 131,
		"LTU": 132,
		"LUX": 133,
		"LVA": 134,
		"MAC": 135,
		"MAF": 136,
		"MAR": 137,
		"MCO": 138,
		"MDA": 139,
		"MDG": 140,
		"MDV": 141,
		"MEX": 142,
		"MHL": 143,
		"MKD": 144,
		"MLI": 145,
		"MLT": 146,
		"MMR": 147,
		"MNE": 148,
		"MNG": 149,
		"MNP": 150,
		"MOZ": 151,
		"MRT": 152,
		"MSR": 153,
		"MTQ": 154,
		"MUS": 155,
		"MWI": 156,
		"MYS": 157,
		"MYT": 158,
		"NAM": 159,
		"NCL": 160,
		"NER": 161,
		"NFK": 162,
		"NGA": 163,
		"NIC": 164,
		"NOR": 165,
		"NIU": 166,
		"NLD": 167,
		"NPL": 168,
		"NRU": 169,
		"NZL": 170,
		"OMN": 171,
		"PAK": 172,
		"PAN": 173,
		"PCN": 174,
		"PER": 175,
		"PHL": 176,
		"PLW": 177,
		"PNG": 178,
		"POL": 179,
		"PRI": 180,
		"PRK": 181,
		"PRT": 182,
		"PRY": 183,
		"PSE": 184,
		"PYF": 185,
		"QAT": 186,
		"REU": 187,
		"ROU": 188,
		"RUS": 189,
		"RWA": 190,
		"SAU": 191,
		"SDN": 192,
		"SEN": 193,
		"SGP": 194,
		"SGS": 195,
		"SHN": 196,
		"SJM": 197,
		"SLB": 198,
		"SLE": 199,
		"SLV": 200,
		"SMR": 201,
		"SOM": 202,
		"SPM": 203,
		"SRB": 204,
		"STP": 205,
		"SUR": 206,
		"SVK": 207,
		"SVN": 208,
		"SWE": 209,
		"SWZ": 210,
		"SYC": 211,
		"SYR": 212,
		"TCA": 213,
		"TCD": 214,
		"TGO": 215,
		"THA": 216,
		"TJK": 217,
		"TKL": 218,
		"TKM": 219,
		"TLS": 220,
		"TON": 221,
		"TTO": 222,
		"TUN": 223,
		"TUR": 224,
		"TUV": 225,
		"TWN": 226,
		"TZA": 227,
		"UGA": 228,
		"UKR": 229,
		"UMI": 230,
		"URY": 231,
		"USA": 232,
		"UZB": 233,
		"VAT": 234,
		"VCT": 235,
		"VEN": 236,
		"VGB": 237,
		"VIR": 238,
		"VNM": 239,
		"VUT": 240,
		"WLF": 241,
		"WSM": 242,
		"YEM": 243,
		"ZAF": 244,
		"ZMB": 245,
		"ZWE": 246,
	}
	cc_country = map[CountryCode]string{
		0:   "UNKOWN",
		1:   "Aruba",
		2:   "Afghanistan",
		3:   "Angola",
		4:   "Anguilla",
		5:   "Åland Islands",
		6:   "Albania",
		7:   "Andorra",
		8:   "Netherlands Antilles",
		9:   "United Arab Emirates",
		10:  "Argentina",
		11:  "Armenia",
		12:  "American Samoa",
		13:  "Antarctica",
		14:  "French Southern Territories",
		15:  "Antigua and Barbuda",
		16:  "Australia",
		17:  "Austria",
		18:  "Azerbaijan",
		19:  "Burundi",
		20:  "Belgium",
		21:  "Benin",
		22:  "Burkina Faso",
		23:  "Bangladesh",
		24:  "Bulgaria",
		25:  "Bahrain",
		26:  "Bahamas",
		27:  "Bosnia and Herzegovina",
		28:  "Saint Barthélemy",
		29:  "Belarus",
		30:  "Belize",
		31:  "Bermuda",
		32:  "Bolivia",
		33:  "Brazil",
		34:  "Barbados",
		35:  "Brunei Darussalam",
		36:  "Bhutan",
		37:  "Bouvet Island",
		38:  "Botswana",
		39:  "Central African Republic",
		40:  "Canada",
		41:  "Cocos (Keeling) Islands",
		42:  "Switzerland",
		43:  "Chile",
		44:  "China",
		45:  "Côte d'Ivoire",
		46:  "Cameroon",
		47:  "Congo, the Democratic Republic of the",
		48:  "Congo",
		49:  "Cook Islands",
		50:  "Colombia",
		51:  "Comoros",
		52:  "Cape Verde",
		53:  "Costa Rica",
		54:  "Cuba",
		55:  "Christmas Island",
		56:  "Cayman Islands",
		57:  "Cyprus",
		58:  "Czechia",
		59:  "Germany",
		60:  "Djibouti",
		61:  "Dominica",
		62:  "Denmark",
		63:  "Dominican Republic",
		64:  "Algeria",
		65:  "Ecuador",
		66:  "Egypt",
		67:  "Eritrea",
		68:  "Western Sahara",
		69:  "Spain",
		70:  "Estonia",
		71:  "Ethiopia",
		72:  "Finland",
		73:  "Fiji",
		74:  "Falkland Islands (Malvinas)",
		75:  "France",
		76:  "Faroe Islands",
		77:  "Micronesia, Federated States of",
		78:  "Gabon",
		79:  "United Kingdom",
		80:  "Georgia",
		81:  "Guernsey",
		82:  "Ghana",
		83:  "Guinea",
		84:  "Gibraltar",
		85:  "Guadeloupe",
		86:  "Gambia",
		87:  "Guinea-Bissau",
		88:  "Equatorial Guinea",
		89:  "Greece",
		90:  "Grenada",
		91:  "Greenland",
		92:  "Guatemala",
		93:  "French Guiana",
		94:  "Guam",
		95:  "Guyana",
		96:  "Hong Kong",
		97:  "Heard Island and McDonald Islands",
		98:  "Honduras",
		99:  "Croatia",
		100: "Haiti",
		101: "Hungary",
		102: "Indonesia",
		103: "Isle of Man",
		104: "India",
		105: "British Indian Ocean Territory",
		106: "Ireland",
		107: "Iran, Islamic Republic of",
		108: "Iraq",
		109: "Iceland",
		110: "Israel",
		111: "Italy",
		112: "Jamaica",
		113: "Jersey",
		114: "Jordan",
		115: "Japan",
		116: "Kazakhstan",
		117: "Kenya",
		118: "Kyrgyzstan",
		119: "Cambodia",
		120: "Kiribati",
		121: "Saint Kitts and Nevis",
		122: "Korea, Republic of",
		123: "Kuwait",
		124: "Lao People's Democratic Republic",
		125: "Lebanon",
		126: "Liberia",
		127: "Libyan Arab Jamahiriya",
		128: "Saint Lucia",
		129: "Liechtenstein",
		130: "Sri Lanka",
		131: "Lesotho",
		132: "Lithuania",
		133: "Luxembourg",
		134: "Latvia",
		135: "Macao",
		136: "Saint Martin (French part)",
		137: "Morocco",
		138: "Monaco",
		139: "Moldova, Republic of",
		140: "Madagascar",
		141: "Maldives",
		142: "Mexico",
		143: "Marshall Islands",
		144: "Macedonia, the former Yugoslav Republic of",
		145: "Mali",
		146: "Malta",
		147: "Myanmar",
		148: "Montenegro",
		149: "Mongolia",
		150: "Northern Mariana Islands",
		151: "Mozambique",
		152: "Mauritania",
		153: "Montserrat",
		154: "Martinique",
		155: "Mauritius",
		156: "Malawi",
		157: "Malaysia",
		158: "Mayotte",
		159: "Namibia",
		160: "New Caledonia",
		161: "Niger",
		162: "Norfolk Island",
		163: "Nigeria",
		164: "Nicaragua",
		165: "Norway",
		166: "Niue",
		167: "Netherlands",
		168: "Nepal",
		169: "Nauru",
		170: "New Zealand",
		171: "Oman",
		172: "Pakistan",
		173: "Panama",
		174: "Pitcairn",
		175: "Peru",
		176: "Philippines",
		177: "Palau",
		178: "Papua New Guinea",
		179: "Poland",
		180: "Puerto Rico",
		181: "Korea, Democratic People's Republic of",
		182: "Portugal",
		183: "Paraguay",
		184: "Palestinian Territory, Occupied",
		185: "French Polynesia",
		186: "Qatar",
		187: "Réunion",
		188: "Romania",
		189: "Russian Federation",
		190: "Rwanda",
		191: "Saudi Arabia",
		192: "Sudan",
		193: "Senegal",
		194: "Singapore",
		195: "South Georgia and the South Sandwich Islands",
		196: "Saint Helena",
		197: "Svalbard and Jan Mayen",
		198: "Solomon Islands",
		199: "Sierra Leone",
		200: "El Salvador",
		201: "San Marino",
		202: "Somalia",
		203: "Saint Pierre and Miquelon",
		204: "Serbia",
		205: "São Tomé and Príncipe",
		206: "Suriname",
		207: "Slovakia",
		208: "Slovenia",
		209: "Sweden",
		210: "Swaziland",
		211: "Seychelles",
		212: "Syrian Arab Republic",
		213: "Turks and Caicos Islands",
		214: "Chad",
		215: "Togo",
		216: "Thailand",
		217: "Tajikistan",
		218: "Tokelau",
		219: "Turkmenistan",
		220: "Timor-Leste",
		221: "Tonga",
		222: "Trinidad and Tobago",
		223: "Tunisia",
		224: "Turkey",
		225: "Tuvalu",
		226: "Taiwan, Province of China",
		227: "Tanzania, United Republic of",
		228: "Uganda",
		229: "Ukraine",
		230: "United States Minor Outlying Islands",
		231: "Uruguay",
		232: "United States",
		233: "Uzbekistan",
		234: "Vatican City State",
		235: "Saint Vincent and the Grenadines",
		236: "Venezuela",
		237: "Virgin Islands, British",
		238: "Virgin Islands, U.S.",
		239: "Viet Nam",
		240: "Vanuatu",
		241: "Wallis and Futuna",
		242: "Samoa",
		243: "Yemen",
		244: "South Africa",
		245: "Zambia",
		246: "Zimbabwe",
	}
)
