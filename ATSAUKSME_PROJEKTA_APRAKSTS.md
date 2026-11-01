# PROJEKTA APRAKSTS UN ATSAUKSME
## Par izstrādāto elektronisko resursu – interaktīvu digitālo spēļu un viktorīnu platformu

---

### 1. Vispārīgā informācija par projektu un līgumu

| Parametrs | Informācija |
| :--- | :--- |
| **Projekta / sistēmas nosaukums** | **"MiniMani" – Reāllaika daudzspēlētāju interaktīvā spēļu un viktorīnu platforma** |
| **Publiski pieejamais resurss** | [https://minimani.kilograms.lv](https://minimani.kilograms.lv) / [https://ugo.lv](https://ugo.lv) |
| **Izstrādātājs / Izpildītājs** | Agris Markus |
| **Līguma / projekta priekšmets** | Interaktīvas, spēļu elementos balstītas tīmekļa viktorīnu un daudzspēlētāju izzinošas platformas projektēšana, programmēšana, mākoņinfrastruktūras izveide un uzturēšana |
| **Izstrādes un ieviešanas periods** | 2026. gads |
| **Projekta statuss** | Pabeigts, veiksmīgi nodots ekspluatācijā, aktīvs un pieejams tiešsaistē |

---

### 2. Atbilstība iepirkuma nolikuma līdzvērtības prasībām

> **Nolikuma prasība:**  
> *"Par līdzvērtīgiem uzskatāmi līgumi, kuru rezultātā izstrādāts elektronisks resurss – interaktīvas digitālās spēles vai viktorīnas, vai spēļu bāzētas platformas, vai interaktīvi tīmekļa risinājumi, vai interaktīvas digitālās izglītības platformas."*

Izstrādātais risinājums **"MiniMani"** pilnībā atbilst visām nolikumā minētajām pazīmēm un ietver katru no definētajiem elementiem:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            MINIMANI PLATFORMA                               │
├──────────────────────┬──────────────────────┬───────────────────────────────┤
│ Interaktīva digitālā │ Spēļu bāzēta         │ Interaktīvs tīmekļa           │
│ spēle un viktorīna   │ platforma (Lobby)    │ risinājums & izziņa           │
├──────────────────────┼──────────────────────┼───────────────────────────────┤
│ • Reāllaika WebSockets│ • Istabu veidošana   │ • Adaptīvs SPA interfeiss     │
│ • 4 Cenu viktorīna   │ • Demokrātiska balss │ • Reālu tirgus datu analīze   │
│ • Precīzā ievade     │ • Spēlētāju pjedestāls│ • Dinamiska punktu sistēma    │
│ • Laika atskaite     │ • Profilu pārvaldība │ • Mākoņglabātuve (S3 WebP)    │
└──────────────────────┴──────────────────────┴───────────────────────────────┘
```

---

### 3. Detalizēts izstrādātā elektroniskā resursa funkcionālais apraksts

#### 3.1. Interaktīva digitālā spēle un viktorīna (Interactive Digital Game & Quiz)
- **Reāllaika sinhronizācija ar zemu aizturi (<50 ms):** Izmantojot `WebSocket` divvirzienu komunikācijas protokolu, spēles gaita, dalībnieku pievienošanās, minējumi un rezultātu atklāšana notiek sinhroni visiem dalībniekiem bez lapas pārlādes.
- **Daudzveidīgi viktorīnas režīmi:**
  - **"4 Cenu Varianti" (Multiple Choice viktorīna):** Algoritmiski tiek atlasīts pareizais sludinājuma parametru un cenas komplekts un dinamiski ģenerēti trīs ticami, matemātiski noapaļoti maldinošie varianti (*distractors*).
  - **"Precīzā cenas ievade":** Spēlētāji paši ar klaviatūru ievada summu eiro, saņemot punktus pēc simetriskās novirzes formulas un papildu "Bullseye" bonusu par precīzu trāpījumu.
- **Interaktīvs raunda taimeris un atgriezeniskā saite:** Katrā raundā darbojas vizuāls dinamiskais progresa taimeris, kas aprēķina papildu ātruma bonusus atkarībā no atbildēšanas sekundes.
- **Raundu atklāšana (Reveal Screen):** Raundam noslēdzoties, lietotājiem interaktīvā grafikā tiek parādīts pareizais rezultāts, visu dalībnieku iesniegtās atbildes un novirze.

#### 3.2. Spēļu bāzēta platforma (Gamified Multiplayer Platform)
- **Privātās un publiskās spēļu istabas (Game Rooms):**
  - Automātiski ģenerēti unikāli 4 burtu/ciparu kodi (piemēram, `AUTO`, `CENA`, `ZIVE`), kas ļauj uzaicināt draugus, kolēģus vai skolēnus ar 1 klikšķi vai QR kodu.
  - Publiskais vestibils un **"Ātrā spēle" (Quick Join)** ērtai tūlītējai dalībai atvērtajās spēlēs.
- **Demokrātiskā vestibila balsošanas mehānika:**
  - Pirms spēles uzsākšanas visi telpā esošie dalībnieki reāllaikā balso par raundu skaitu (5, 10, 15, 20) un spēles režīmu.
  - Tiklīdz kāds variants iegūst $>50\%$ balsu, tas automātiski kļūst par apstiprināto istabas noteikumu.
- **Gaitas un uzvarētāju pjedestāls (Podium & Leaderboard):**
  - Noslēguma ekrānā tiek ģenerēts 1., 2. un 3. vietas pjedestāls ar animētiem avatariem un konfeti vizuālajiem efektiem.
  - Pilna visu dalībnieku rezultātu kopsavilkuma tabula ar punktu detalizāciju.
- **Spēlētāju profili un autorizācija:**
  - Hibrīda piekļuves modelis: iespēja spēlēt gan tūlītēji kā viesim (*guest access*), gan reģistrēt pastāvīgu lietotāja profilu ar e-pastu, publisko niku, unikālu avatara ikonu un spēļu vēstures statistiku.

#### 3.3. Interaktīvs tīmekļa risinājums un izglītības/izziņas komponente
- **Vērtēšanas un tirgus izpratnes digitālā izglītība:** Spēle kalpo kā interaktīvs izziņas instruments, trenējot dalībnieku spēju analizēt nekustamo īpašumu un materiālo vērtību parametrus (atrašanās vieta, ēkas platība, zemes platība, stāvu skaits, istabu skaits, inženierkomunikācijas un labiekārtojums).
- **Informatīvā noteikumu un formulas lapa:** Platformā integrēts interaktīvs noteikumu modulis ([RulesModal](https://minimani.kilograms.lv/api/rules)), kas tiešsaistē atspoguļo matemātisko algoritmu un punktu aprēķina loģiku ar praktiskiem piemēriem.
- **Mobilajām ierīcēm pielāgota pieredze (Responsive Design):** Risinājums ir 100% pielāgots viedtālruņiem, planšetdatoriem un skārienjutīgajiem ekrāniem, neprasot lietotājiem lejupielādēt un instalēt vietējās mobilās lietotnes.

---

### 4. Tehnoloģiskā arhitektūra un ieviestie standarti

Projekta izstrādē pielietotas mūsdienīgas, industrijā atzītas atvērtā koda tehnoloģijas un drošības standarti:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        KLIENTA PĀRLŪKPROGRAMMA                         │
│             React 18 · TypeScript · Vite · Tailwind CSS                │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ HTTPS (TLS 1.3) / WSS
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                   REVERSAIS PROXY UN TĪKLA SLĀNIS                      │
│                  Nginx (HTTP/2, Gzip, SSL, SPA maršrutēšana)           │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ 127.0.0.1:8088
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                        BACKEND SERVERIS (Go)                           │
│     Chi REST API · In-Memory WebSocket Hub · S3 Pure-Go Uploader       │
└───────────────────┬────────────────────────────────┬───────────────────┘
                    │                                │
                    ▼                                ▼
┌──────────────────────────────────────┐  ┌──────────────────────────────┐
│       POSTGRESQL 16 & POSTGIS        │  │   HETZNER OBJECT STORAGE     │
│   Sludinājumi, profili, rezultāti    │  │   WebP attēli (RAM straume)  │
└──────────────────────────────────────┘  └──────────────────────────────┘
```

1. **Lietotāja saskarne (Frontend):**
   - **React 18** un **TypeScript** nodrošina stingru tipu drošību un ātru renderingu.
   - **Tailwind CSS** garantē adaptīvu (responsive) izskatu tumšajā dizaina tēmā (*dark mode*).
   - **Canvas Confetti** un **Lucide Icons** vizuālās iesaistes un geimifikācijas dinamikai.
2. **Pakalpojumu serveris (Backend):**
   - Izstrādāts **Go (Golang)** programmēšanas valodā, garantējot augstu izpildes ātrumu, minimālu RAM patēriņu un tūkstošiem paralēlo WebSocket savienojumu apkalpošanu.
   - Droša paroļu šifrēšana ar **bcrypt** un autentifikācija ar **JWT** (JSON Web Tokens).
3. **Datu un failu glabātuve:**
   - **PostgreSQL 16** relāciju datubāze ar **PostGIS 3.4** telpisko datu paplašinājumu.
   - **Hetzner Object Storage (S3 saderīga):** Attēlu apstrāde un konvertācija uz **WebP** formātu notiek tikai operatīvajā atmiņā (zero disk usage), novēršot servera diska noslodzi un nodrošinot augstu ielādes ātrumu.
4. **DevOps un izvietošana (CI/CD):**
   - **Docker Compose** konteinerizācija neatkarīgai vides darbībai.
   - Automatizēts **GitHub Actions** piegādes cauruļvads (CI/CD workflow), kas pie koda atjaunināšanas automātiski veic kompilāciju, testus un izvietošanu uz produkcijas servera.
   - **Let's Encrypt SSL/TLS** šifrēšana ar automātisku sertifikātu atjaunošanu.

---

### 5. Sasniegtie rezultāti un veiktspējas rādītāji

- **Tūlītējs spēles starts:** Datu bāzes kataloga pirmrindas arhitektūra nodrošina spēles sesijas uzsākšanu mazāk nekā 100 milisekundēs.
- **Daudzspēlētāju mērogojamība:** Platforma spēj paralēli uzturēt desmitiem vienlaicīgu istabu ar 8+ spēlētājiem katrā.
- **Pieejamība un stabilitāte:** Nodrošināts augsts pakalpojuma pieejamības līmenis ($99.9\%$), zems servera resursu patēriņš un pilnīga datu integritāte.
- **Lietotāju atsauksmes:** Pozitīvs novērtējums par spēles dinamiku, vizuālo pievilcību, intuitīvo saskarni un izzinošo vērtību.

---

### 6. Apliecinājums / Pasūtītāja atsauksme

Ar šo tiek apliecināts, ka:

1. Izpildītājs **Agris Markus** ir pilnā apjomā, kvalitatīvi un līgumā noteiktajā termiņā izstrādājis, notestējis un ieviesis ekspluatācijā elektronisko resursu – interaktīvu digitālo daudzspēlētāju spēļu un viktorīnu platformu **"MiniMani"**.
2. Izstrādātā sistēma darbojas stabili, bez kļūdām un atbilst visām augstākā līmeņa tehniskajām, drošības un interaktivitātes prasībām.
3. Pretenziju par veiktajiem darbiem, to kvalitāti vai izpildes termiņiem nav.
4. Izstrādātais risinājums pilnībā atbilst iepirkuma nolikumā definētajai līdzvērtīga līguma definīcijai: **elektronisks resurss – interaktīva digitālā spēle / viktorīna / spēļu bāzēta platforma / interaktīvs tīmekļa risinājums**.

---

**Atsauksmes sagatavošanas datums:** 2026. gada 1. novembrī  
**Projekta tiešsaistes adrese:** [https://minimani.kilograms.lv](https://minimani.kilograms.lv) / [https://ugo.lv](https://ugo.lv)
