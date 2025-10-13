package main

import (
	"3xui-bot/internal/adapters/marzban"
	"3xui-bot/internal/core"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	repo := marzban.NewMarzbanRepository("https://carrot-promo.ru", "3vil3vil3vil", "BewareOfDoggg1422")

	if err := repo.Authenticate(ctx); err != nil {
		log.Fatal(err)
	}

	inbounds, err := repo.GetInbounds(ctx)
	if err != nil {
		log.Fatal("Failed to get inbounds:", err)
	}

	fmt.Printf("\n📋 Found %d inbounds:\n", len(inbounds))
	for i, inbound := range inbounds {
		inboundJSON, _ := json.MarshalIndent(inbound, "  ", "  ")
		fmt.Printf("\n%d. %s\n", i+1, string(inboundJSON))

		// Выводим ключевые поля
		if tag, ok := inbound["tag"].(string); ok {
			fmt.Printf("   Tag: %s\n", tag)
		}
		if protocol, ok := inbound["protocol"].(string); ok {
			fmt.Printf("   Protocol: %s\n", protocol)
		}
	}

	// Если нет inbounds, используем пустой массив (возможно сервер разрешает любые)
	if len(inbounds) == 0 {
		fmt.Println("\n⚠️  No specific inbounds found. Trying without inbound restriction...")
	}

	// Извлекаем теги inbounds для создания пользователя
	var inboundTags []string
	inboundsByProtocol := make(map[string][]string)

	for _, inbound := range inbounds {
		if tag, ok := inbound["tag"].(string); ok {
			inboundTags = append(inboundTags, tag)
			if protocol, ok := inbound["protocol"].(string); ok {
				inboundsByProtocol[protocol] = append(inboundsByProtocol[protocol], tag)
			}
		}
	}

	if len(inboundTags) > 0 {
		fmt.Println("\n📝 Available inbound tags:", inboundTags)
		fmt.Println("📝 Inbounds by protocol:")
		for proto, tags := range inboundsByProtocol {
			fmt.Printf("   %s: %v\n", proto, tags)
		}
	}

	// Подготовка данных пользователя
	dataLimit := int64(10 * 1024 * 1024 * 1024)            // 10 GB
	expireTimestamp := time.Now().AddDate(0, 0, 30).Unix() // 30 дней

	// Создаем структуру inbounds используя реальные теги
	userInbounds := make(map[string][]string)
	for protocol, tags := range inboundsByProtocol {
		userInbounds[protocol] = tags
	}

	// Если нет inbounds, пробуем создать пользователя без них или с пустым объектом
	if len(userInbounds) == 0 {
		fmt.Println("\n⚠️  No inbounds configured, creating user without specific inbounds")
		// Можно попробовать пустой объект или nil
	}

	user := core.MarzbanUserData{
		Username:  "test_user_" + fmt.Sprintf("%d", time.Now().Unix()),
		DataLimit: &dataLimit,
		Expire:    &expireTimestamp,
		Status:    "active",
		Note:      "Test user from Go client",
		Proxies: map[string]interface{}{
			"vless": map[string]interface{}{},
		},
		Inbounds: userInbounds,
	}

	// Логируем что отправляем
	userJSON, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println("\n📤 Sending user data:")
	fmt.Println(string(userJSON))

	fmt.Println("\n🔄 Creating user...")
	createdUser, err := repo.CreateUser(ctx, &user)
	if err != nil {
		log.Fatal("❌ Create user failed:", err)
	}

	// Логируем ответ
	responseJSON, _ := json.MarshalIndent(createdUser, "", "  ")
	fmt.Println("\n✅ User created successfully!")
	fmt.Println("📥 Response:")
	fmt.Println(string(responseJSON))

	// Выводим ссылки для подключения
	if len(createdUser.Links) > 0 {
		fmt.Println("\n🔗 Connection links:")
		for i, link := range createdUser.Links {
			fmt.Printf("%d. %s\n", i+1, link)
		}
	}

	if createdUser.SubscriptionURL != "" {
		fmt.Printf("\n📱 Subscription URL: %s\n", createdUser.SubscriptionURL)
	}

	fmt.Println("\n🔍 Testing GetUser...")
	fetchedUser, err := repo.GetUser(ctx, createdUser.Username)
	if err != nil {
		log.Fatal("❌ Get user failed:", err)
	}

	fmt.Println("✅ User fetched successfully!")
	fmt.Printf("   Username: %s\n", fetchedUser.Username)
	fmt.Printf("   Status: %s\n", fetchedUser.Status)
	if fetchedUser.DataLimit != nil {
		fmt.Printf("   Data Limit: %d GB\n", *fetchedUser.DataLimit/(1024*1024*1024))
	}
}
