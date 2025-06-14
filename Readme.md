# 💳 Midtrans Integration Project

Proyek ini merupakan implementasi backend integrasi **Midtrans Payment Gateway** dengan dua jenis pembayaran utama:

1. **Pembayaran Biasa (One-Time Payment)**: Cocok untuk produk/jasa yang dibayar sekali.
2. **Pembayaran Berkala (Subscription/Recurring Payment)**: Cocok untuk layanan berlangganan.

Proyek ini menangani transaksi, notifikasi (webhook), refund, dan pembatalan subscription.

---

## 📦 Fitur

### 🔹 One-Time Payment
- Membuat transaksi
- Menerima notifikasi pembayaran dari Midtrans
- Melakukan refund berdasarkan `order_id`

### 🔸 Subscription
- Membuat subscription billing 
- Cancel subscription
- **Catatan**: Midtrans tidak mengirim webhook otomatis untuk subscription. Harus **manual cek status** atau gunakan cron polling.

---

## ⚙️ Setup Awal

### 1. Login ke Midtrans & Dapatkan Server Key
1. Buka [https://dashboard.midtrans.com/](https://dashboard.midtrans.com/)
2. Masuk ke akun kamu
3. Pilih menu: `Settings > Access Keys`
4. Salin **Server Key**
5. Simpan ke environment:
   ```env
   MIDTRANS_SERVER_KEY=Your_Server_Key
   MIDTRANS_CLIENT_KEY=Your_Client_Key (opsional)
   MIDTRANS_ENV=sandbox
