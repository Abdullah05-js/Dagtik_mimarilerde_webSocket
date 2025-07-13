🧠 Dağıtık Gerçek Zamanlı WebSocket Chat Uygulaması
Bu proje, Go (Golang) dili kullanılarak geliştirilen dağıtık mimariye sahip bir gerçek zamanlı sohbet (chat) sistemidir. WebSocket, Apache Kafka ve MongoDB kullanılarak, yatayda ölçeklenebilir, performanslı ve güvenilir bir chat altyapısı sunar.

🚀 Özellikler
⚡ Gerçek zamanlı mesajlaşma (WebSocket üzerinden)

🔄 Kafka ile dağıtık mesaj işleme ve yük dengeleme

🧩 MongoDB ile mesajların kalıcı olarak saklanması

🌐 Fiber framework ile hızlı ve modern bir sunucu

🧵 Goroutine, Mutex ve Context ile yüksek eşzamanlılık

♻️ Çoklu sunucu örneğiyle (instance) çalışma desteği

🔐 (Opsiyonel) WebSocket bağlantısında JWT kimlik doğrulama

🪝 Kafka consumer ve producer yapıları

💚 react-use-websocket ile güçlü ve sade bir frontend

🏗️ Mimarisi
[Client] ⇄ [Fiber + WebSocket Sunucusu] ⇄ [Kafka Topic] ⇄ [Diğer Sunucular] ⇄ [MongoDB]

Kullanıcı, WebSocket ile bağlanır ve mesaj gönderir.

Mesaj Kafka topic’e yazılır.

Tüm backend sunucuları bu topic'i dinler.

Mesajlar hem MongoDB’ye kaydedilir hem de ilgili kullanıcılara iletilir.

🛠️ Kullanılan Teknolojiler
Teknoloji

Açıklama

Go (Golang)

Backend dili

Fiber

HTTP & WebSocket sunucusu

Kafka

Dağıtık mesaj kuyruğu

MongoDB

Veritabanı

WebSocket

Gerçek zamanlı bağlantı

JWT (opsiyonel)

📁 Proje Yapısı
ChatAPP/
├── kafka/ # Kafka reader, writer ve client bağlantıları
├── model/ # Veri şemaları (MessageSchema vb.)
├── database/ # MongoDB bağlantısı
├── main.go # Fiber sunucusunun başlangıç noktası
├── go.mod # Go modül yönetimi

📦 Kurulum
🧰 Gerekli Araçlar
Go

Kafka & Zookeeper

MongoDB

🔧 Kafka’yı Elle Kurmak (Docker Yoksa)
Kafka ve Zookeeper'ı indirin: https://kafka.apache.org/downloads

Aşağıdaki komutları sırayla farklı terminallerde çalıştırarak servisleri başlatın:

# 1. Terminal: Zookeeper'ı başlat

bin/zookeeper-server-start.sh config/zookeeper.properties

````bash
# 2. Terminal: Kafka'yı başlat
bin/kafka-server-start.sh config/server.properties
```bash
# 3. Terminal: Gerekli topic'i oluştur
bin/kafka-topics.sh --create --topic chat-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

▶️ Çalıştırma
Backend
# Bağımlılıkları yükle
go mod tidy

# Sunucuyu başlat
go run server.go

🌍 API Uç Noktaları
Yöntem

URL

Açıklama

GET

/ws

WebSocket bağlantı noktası. Kimlik doğrulama için query parametresi olarak token (JWT) gönderilebilir.

🛡️ Güvenlik
Kullanıcı kimliğini doğrulamak için WebSocket bağlantısı sırasında JWT kullanılabilir (opsiyonel).

Bir istemcinin WebSocket bağlantısı kesildiğinde, sunucudaki ilgili kaydı güvenli bir şekilde silinir.

Geçersiz formatta gönderilen mesajlar işlenmez ve kullanıcıya hata mesajı döndürülür.

🧪 Geliştirici Notları
Aktif istemciler, map[bson.ObjectID]ConnectionInfo yapısı kullanılarak sunucu belleğinde takip edilir.

Tüm sunucu örnekleri (instance) arasındaki mesaj senkronizasyonu, doğrudan bir bağlantı yerine Kafka üzerinden sağlanır. Bu, sistemin dağıtık yapısını korur.

💡 Katkıda Bulunma
Katkılarınızı memnuniyetle karşılıyoruz! Projeyi geliştirmek için hata bildirimi (issue) açabilir veya yeni özellikler/düzeltmeler için Pull Request (PR) gönderebilirsiniz.

📜 Lisans
Bu proje MIT Lisansı ile lisanslanmıştır.

✨ Geliştirici
👨‍💻 Abdullah Han

📧 AbdullahHan05@proton.me
````
