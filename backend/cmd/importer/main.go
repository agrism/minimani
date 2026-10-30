package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	internaldb "minimani/backend/internal/db"
	db "minimani/backend/internal/db/sqlc"
	"minimani/backend/internal/ss"
)

type SampleListing struct {
	Category    string
	Title       string
	Description string
	Price       float64
	Region      string
	ImageURL    string
	OriginalURL string
}

func main() {
	_ = godotenv.Load("../.env", ".env")

	action := flag.String("action", "seed", "Action to run: 'seed' or 'scrape'")
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://minimani:minimani_secret@localhost:5438/minimani_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Neizdevās pieslēgties datubāzei: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	if *action == "seed" {
		seedDatabase(ctx, queries)
	} else {
		fmt.Printf("Nezināma darbība: %s\n", *action)
	}
}

func seedDatabase(ctx context.Context, queries *db.Queries) {
	log.Println("🌱 Sākam ss.lv parauga sludinājumu sēšanu...")

	samples := getCuratedListings()

	inserted := 0
	for _, s := range samples {
		optionsJSON := ss.GenerateOptionsJSON(s.Price)

		_, err := queries.InsertListing(ctx, db.InsertListingParams{
			Category:    s.Category,
			Title:       s.Title,
			Description: pgtype.Text{String: s.Description, Valid: s.Description != ""},
			ImagePath:   s.ImageURL,
			OriginalUrl: pgtype.Text{String: s.OriginalURL, Valid: s.OriginalURL != ""},
			ActualPrice: internaldb.NumericFromFloat(s.Price),
			Currency:    "EUR",
			Region:      pgtype.Text{String: s.Region, Valid: s.Region != ""},
			Options:     optionsJSON,
		})
		if err != nil {
			log.Printf("Kļūda ievietojot '%s': %v", s.Title, err)
			continue
		}
		inserted++
	}

	total, _ := queries.CountListings(ctx)
	log.Printf("✅ Veiksmīgi pievienoti %d sludinājumi! Kopā aktīvi: %d\n", inserted, total)
}

func getCuratedListings() []SampleListing {
	return []SampleListing{
		// ==================== AUTO ====================
		{
			Category:    "auto",
			Title:       "BMW 330d E46 M-Pakotne, 2003. gads",
			Description: "Klasisks 3.0 dīzelis 150kw. Melns ādas salons, melnie griesti, R18 diski ar labām riepām. Tikko izieta TA bez aizrādījumiem. Sliekšņi pārmetināti, čipots nav.",
			Price:       3850,
			Region:      "Rīga",
			ImageURL:    "https://images.unsplash.com/photo-1555215695-3004980ad54e?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/cars/bmw/330/index.html",
		},
		{
			Category:    "auto",
			Title:       "Volkswagen Golf 4, 1.9 TDI (81 kW), 1999. gads",
			Description: "Leģendārais sūkņa dīzelis. Ļoti ekonomisks (4.5l/100km). Ir vizuāli defekti uz priekšējiem spārniem, bet tehniskais stāvoklis lielisks. Ielikts jauns sajūga komplekts.",
			Price:       1350,
			Region:      "Jelgava",
			ImageURL:    "https://images.unsplash.com/photo-1541899481282-d53bffe3c35d?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/cars/volkswagen/golf-4/index.html",
		},
		{
			Category:    "auto",
			Title:       "Audi A6 C6 3.0 TDI Quattro S-Line, 2006. gads",
			Description: "Pilnpiedziņa strādā nevainojami. Pneimo piekare pārbaudīta, sprauslas nomainītas pirms gada. BOSE audio sistēma, salons kopts, oriģināls nobraukums 340 000 km.",
			Price:       4900,
			Region:      "Valmiera",
			ImageURL:    "https://images.unsplash.com/photo-1603584173870-7f23fdae1b7a?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/cars/audi/a6/index.html",
		},
		{
			Category:    "auto",
			Title:       "Tesla Model 3 Long Range AWD, 2021. gads",
			Description: "Baterijas atlikums 94%. Baltā pērļu krāsa, Premium melnais salons. Autopilots, siltumsūknis, apsildāmi visi sēdekļi. Eiropas modelis, servisa vēsture pieejama.",
			Price:       27900,
			Region:      "Mārupe",
			ImageURL:    "https://images.unsplash.com/photo-1560958089-b8a1929cea89?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/cars/tesla/model-3/index.html",
		},
		{
			Category:    "auto",
			Title:       "Porsche Cayenne S 4.5 V8, 2004. gads",
			Description: "Jaudīgs V8 motors (340 ZS) ar uzstādītu gāzes iekārtu (LPG). Ādas salons, lūka, 20 collu diski. Dodu līdzi ziemas riepu komplektu uz diskiem.",
			Price:       5400,
			Region:      "Liepāja",
			ImageURL:    "https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/cars/porsche/cayenne/index.html",
		},
		{
			Category:    "auto",
			Title:       "VAZ 2101 'Žigulis', 1978. gads (Kolekcijas stāvoklis)",
			Description: "Oriģināls krāsojums, 54 000 km nobraukums. Ziemā nav braukts, glabāts siltā garāžā. Rūpnīcas hroms ideālā stāvoklī. Vēsturiskā spēkrata statuss.",
			Price:       4500,
			Region:      "Cēsis",
			ImageURL:    "https://images.unsplash.com/photo-1583121274602-3e2820c69888?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/cars/vaz/2101/index.html",
		},

		// ==================== NEKUSTAMIE ĪPAŠUMI ====================
		{
			Category:    "real_estate",
			Title:       "2-istabu dzīvoklis Purvciemā, 103. sērija (52 m²)",
			Description: "Izolētas istabas, iestiklota lodžija. Dzīvoklis bez remonta, bet tīrs. Zeme īpašumā, zemi komunālie maksājumi. Blakus skolas, bērnudārzi un sabiedriskais transports.",
			Price:       52000,
			Region:      "Rīga, Purvciems",
			ImageURL:    "https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/real-estate/flats/riga/purvciems/index.html",
		},
		{
			Category:    "real_estate",
			Title:       "Dizaina dzīvoklis Klusajā centrā, renovēta jūgendstila ēka (88 m²)",
			Description: "3 istabas, augstie griesti (3.4m), restaurētas vēsturiskās parketa grīdas un podiņu krāsns. Autonoma gāzes apkure, lifts, slēgts pagalms ar autostāvvietu.",
			Price:       245000,
			Region:      "Rīga, Centrs",
			ImageURL:    "https://images.unsplash.com/photo-1600585154340-be6161a56a0c?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/real-estate/flats/riga/centre/index.html",
		},
		{
			Category:    "real_estate",
			Title:       "Jauna vienstāva privātmāja Mārupē (160 m², zeme 1200 m²)",
			Description: "A+ energoefektivitāte, gaiss-ūdens siltumsūknis, rekuperācijas sistēma. 3 guļamistabas, plaša viesistaba ar izeju uz terasi, garāža 1 auto. Nodota ekspluatācijā 2023. gadā.",
			Price:       285000,
			Region:      "Mārupe",
			ImageURL:    "https://images.unsplash.com/photo-1580587771525-78b9dba3b914?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/real-estate/homes-summer-residences/marupe/index.html",
		},
		{
			Category:    "real_estate",
			Title:       "Lauku viensēta ar pirti un dīķi pie Cēsīm (Zeme 3.5 ha)",
			Description: "Autentiska guļbūves dzīvojamā māja, funkcionējoša melnā pirts, liels zivju dīķis un iekopts ābeļdārzs. Labs piebraucamais ceļš, 3 fāžu elektrība, akas ūdens.",
			Price:       69000,
			Region:      "Cēsu novads",
			ImageURL:    "https://images.unsplash.com/photo-1518780664697-55e3ad937233?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/real-estate/homes-summer-residences/cesis/index.html",
		},

		// ==================== TEHNIKA & BŪVNIECĪBA ====================
		{
			Category:    "tech",
			Title:       "Traktors Belarus MTZ-82 ar frontālo iekrāvēju, 1994. gads",
			Description: "Darba kārtībā, lielā kabīne, jauns starteris un akumulators. Pilnpiedziņa darbojas, hidraulika sausa. TA derīga līdz rudenim. Labs palīgs saimniecībā.",
			Price:       7800,
			Region:      "Kuldīga",
			ImageURL:    "https://images.unsplash.com/photo-1592982537447-7440770cbfc9?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/agriculture/agricultural-machinery/tractors/belarus/index.html",
		},
		{
			Category:    "tech",
			Title:       "Zāles pļāvējs-traktoriņš Husqvarna TS 138",
			Description: "Pļaušanas platums 97 cm, jaudīgs Briggs & Stratton dzinējs ar piespiedu eļļošanu. Lietots 3 sezonas piemājas zālienam (ap 2000 m²). Veikta pilna apkope pirms sezonas.",
			Price:       1850,
			Region:      "Ogre",
			ImageURL:    "https://images.unsplash.com/photo-1589923188900-85dae523342b?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/construction/garden-technics/mowers-riders/index.html",
		},
		{
			Category:    "tech",
			Title:       "iPhone 15 Pro Max 256GB Natural Titanium (Kā jauns)",
			Description: "Baterijas veselība 98%, no pirmās dienas vāciņā un ar aizsargstiklu. Pilns komplekts ar oriģinālo kārbiņu un vadu. Pirkts LMT, oficiālā garantija vēl 14 mēnešus.",
			Price:       990,
			Region:      "Rīga",
			ImageURL:    "https://images.unsplash.com/photo-1695048133142-1a20484d2569?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/electronics/telephones/mobile-phones/apple/iphone-15-pro-max/index.html",
		},

		// ==================== PĒRLES & KURIOZI ====================
		{
			Category:    "pearls",
			Title:       "Autentisks PSRS sienas vilnas tepiķis 'Brieži pie ezera'",
			Description: "Klasika no 80-tajiem gadiem. Stāvējis pie sienas viesistabā, kodes nav ēdušas. Košas krāsas, bārkstis veselas. Ideāls rekvizīts fotosesijām vai retro interjeram.",
			Price:       75,
			Region:      "Daugavpils",
			ImageURL:    "https://images.unsplash.com/photo-1600121848594-d8644e57abab?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/home-stuff/furniture-interior/carpets/index.html",
		},
		{
			Category:    "pearls",
			Title:       "Paštaisīts purva visurgājējs ar VAZ motoru un milzu riepām",
			Description: "Konstruēts medībām un purviem. Peld pa ūdeni! Pazeminātie ātrumi, mehāniskā vinča priekšā. Dokumentu nav un nebūs, bet brauc visur, kur traktors grimst.",
			Price:       2200,
			Region:      "Madona",
			ImageURL:    "https://images.unsplash.com/photo-1542282088-72c9c27ed0cd?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/transport/other/transport-with-defects-or-after-crash/index.html",
		},
		{
			Category:    "pearls",
			Title:       "Restaurēta Ulmaņlaika ozolkoka bufete ar kristāla stikliem",
			Description: "Izgatavota 1930-tajos gados Latvijā. Roku darba kokgriezumi, oriģinālās misiņa furnitūras ar atslēdziņām. Profesionāli restaurēta un vaskota.",
			Price:       1150,
			Region:      "Tukums",
			ImageURL:    "https://images.unsplash.com/photo-1538688525198-9b88f6f53126?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/home-stuff/antique/furniture/index.html",
		},
		{
			Category:    "pearls",
			Title:       "Čuguna malkas plīts 'Zvaigzne' ar cepeškrāsni",
			Description: "Demontēta no lauku virtuves. Čuguna virsma un riņķi bez plaisām, kurtuves durtiņas slēdzas blīvi. Gatava montāžai jaunā vietā.",
			Price:       180,
			Region:      "Bauska",
			ImageURL:    "https://images.unsplash.com/photo-1544816155-12df9643f363?auto=format&fit=crop&w=800&q=80",
			OriginalURL: "https://www.ss.lv/msg/lv/construction/heating-systems/stoves-fireplaces/index.html",
		},
	}
}
