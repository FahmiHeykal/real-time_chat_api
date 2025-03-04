# Real-Time Chat API dengan WebSocket & PostgreSQL

Proyek ini adalah **Real-time Chat API** menggunakan **Golang, WebSocket, dan PostgreSQL**.  
Dibuat sebagai latihan untuk memahami **WebSocket, database PostgreSQL, dan komunikasi real-time** dalam aplikasi backend.

## Fitur
- **WebSocket** untuk komunikasi real-time  
- **Simpan pesan ke PostgreSQL**  
- **Ambil riwayat chat terakhir (50 pesan)**  
- **Mendukung banyak pengguna secara real-time**  

## Teknologi yang Digunakan
- **Golang** (Backend utama)
- **WebSocket** (Komunikasi real-time)
- **PostgreSQL** (Database penyimpanan chat)
- **HTML + JavaScript** (Client sederhana untuk uji coba)

## Cara Menjalankan

###  Persiapan Database PostgreSQL
Pastikan PostgreSQL sudah berjalan, lalu buat database dan tabel:  
```sql
CREATE DATABASE chatdb;
\c chatdb
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW()
);
```

## Clone Proyek & Jalankan Server
```
git clone https://github.com/FahmiHyekal/real-time_chat_api.git
cd real-time_chat_api
```

# Jalankan server
go run main.go


## Uji Coba di Postman atau HTML

Gunakan WebSocket request ke :

```ws://localhost:8080/ws```


# Kirim JSON seperti ini:

```{ "username": "nama anda", "message": "Halo, ini pesan pertama!" }```


# Uji dengan HTML Client

<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <title>Real-Time Chat</title>
</head>
<body>
    <h2>Chat Room</h2>
    <input type="text" id="message" placeholder="Ketik pesan...">
    <button onclick="sendMessage()">Kirim</button>
    <ul id="chat"></ul>
    <script>
        let ws = new WebSocket("ws://localhost:8080/ws");
        ws.onmessage = (event) => {
            let msg = JSON.parse(event.data);
            let li = document.createElement("li");
            li.textContent = msg.username + ": " + msg.message;
            document.getElementById("chat").appendChild(li);
        };
        function sendMessage() {
            let msg = document.getElementById("message").value;
            ws.send(JSON.stringify({ username: "User1", message: msg }));
            document.getElementById("message").value = "";
        }
    </script>
</body>
</html>
