# MiniMani 🏷️

> **"MiniMani"** ir reāllaika daudzspēlētāju cenu minēšanas spēle (*Kahoot / Jackbox* stils) ar sludinājumiem ekskluzīvi no **ss.lv**.

Spēle ir mobilajām ierīcēm pielāgota tīmekļa lietotne latviešu valodā, kurā var spēlēt gan privātās istabās ar draugiem (izmantojot 4-burtu kodu vai saiti), gan pievienoties atvērtajām publiskajām spēlēm ar funkciju **"Ātrā spēle"**.

---

## 🎮 Galvenās Iespējas

- **Istabu veidi:**
  - **Privātās istabas:** Piekļuve ar 4 burtu kodu vai uzaicinājuma saiti / QR kodu.
  - **Atvērtās spēles:** Publisko istabu saraksts sākumlapā ar spēlētāju skaitītājiem.
  - **Ātrā spēle (Quick Play):** 1 klikšķa savienošanās ar pirmo pieejamo istabu.
- **Demokrātiskā vestibila balsošana:**
  - Dalībnieki balso par **jautājumu skaitu (5, 10, 15, 20)** un **režīmu**.
  - Reāllaika balsu sadalījuma procentu joslas. Kad kāds variants sasniedz $>50\%$, tas automātiski tiek apstiprināts.
- **Divi spēles režīmi:**
  - **"Ierakstīt cenu":** Precīza summas ievade ar ātrajām pogām (`+100`, `+1k`, `+5k`, `+10k`). Punkti par precizitāti (0–1000) + *Bullseye* ($\pm 5\%$) un ātruma bonusi.
  - **"4 Atbilžu varianti":** Dinamiska viktorīna ar 4 krāsainām cenu kartītēm (serveris algoritmiski ģenerē ticamas, noapaļotas maldinošās cenas).
- **ss.lv Sludinājumu banka:**
  - Kategorijas: *Auto & Transports*, *Nekustamie īpašumi*, *Tehnika & Būvniecība*, *ss.lv Pērles & Kuriozi*.
  - Attēli tiek arhivēti un kešoti lokāli, novēršot ss.lv hotlink bloķēšanu un 404 kļūdas.
- **Dramatiskā atklāšana un Pjedestāls:**
  - Katra raunda beigās tiek atklāta ss.lv patiesā cena un visu dalībnieku minējumi ar novirzi.
  - Spēles noslēgumā: 1., 2. un 3. vietas pjedestāls ar konfeti salūtu un poga "Spēlēt vēlreiz".

---

## 🛠️ Tehnoloģiskais Steks

- **Frontend:** React 18, Vite, TypeScript, Tailwind CSS, Lucide Icons, Canvas Confetti
- **Backend:** Go (Chi maršrutētājs, `coder/websocket` reāllaika dzinējs)
- **Datubāze:** PostgreSQL 16 ar PostGIS paplašinājumu
- **Kešatmiņa / Ziņojumi:** Redis 7
- **Reversais proxy / HTTPS:** Caddy
- **Konteinerizācija:** Docker Compose (izvietošanai uz Hetzner VPS)

---

## 🚀 Kā Palaist Lokāli

### 1. Klonēt repozitoriju
```bash
git clone https://github.com/agrism/minimani.git
cd minimani
```

### 2. Palaist PostgreSQL & Redis konteinerus
```bash
docker compose -f docker-compose.dev.yml up -d
```

### 3. Ielādēt sākotnējos ss.lv sludinājumus
```bash
cd backend
go run cmd/importer/main.go -action=seed
```

### 4. Palaist Go Backend serveri
```bash
cd backend
go run cmd/server/main.go
# Serveris darbosies uz: http://localhost:8088
```

### 5. Palaist React Frontend
```bash
cd frontend
npm install
npm run dev
# Frontend būs pieejams: http://localhost:5185
```

Atveriet pārlūkprogrammā **http://localhost:5185/** un izbaudiet spēli!

---

## 🌐 Izvietošana Ražošanā (Production ar Docker Compose & Hetzner)

```bash
cp .env.example .env
# Iestatiet DOMAIN un ACME_EMAIL iekš .env
docker compose -f docker-compose.prod.yml up -d --build
```

