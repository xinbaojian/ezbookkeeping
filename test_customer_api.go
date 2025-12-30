// 测试客户管理API的简单脚本
// 运行前先确保应用已启动

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		Id       int64  `json:"id,string"`
		Username string `json:"username"`
	} `json:"user"`
}

type CustomerCreateRequest struct {
	ClientSessionId string `json:"client_session_id"`
	Name            string `json:"name"`
	CustomerType    int8   `json:"customer_type"`
	Address         string `json:"address"`
	Contacts        string `json:"contacts"`
	ContactsInfo    string `json:"contacts_info"`
	Comment         string `json:"comment"`
	Hidden          bool   `json:"hidden"`
}

type CustomerCreateResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	CustomerType int8   `json:"customer_type"`
	Address      string `json:"address"`
	Contacts     string `json:"contacts"`
	ContactsInfo string `json:"contacts_info"`
	Comment      string `json:"comment"`
	Hidden       bool   `json:"hidden"`
	CreatedTime  string `json:"created_time"`
	UpdatedTime  string `json:"updated_time"`
}

func main() {
	// 1. 注册并登录用户
	fmt.Println("=== 测试客户管理API ===")
	fmt.Println("\n1. 注册用户...")

	// 假设直接注册用户
	fmt.Println("请先在Web UI上注册一个用户，或使用以下默认账户登录")
	fmt.Println("用户名: demo 密码: demo")

	// 2. 登录获取token
	fmt.Println("\n2. 登录用户...")
	token := loginUser()
	if token == "" {
		fmt.Println("登录失败，请先注册用户或检查用户名密码")
		return
	}
	fmt.Printf("登录成功，token: %s...\n", token[:20]+"...")

	// 3. 创建客户
	fmt.Println("\n3. 创建客户...")
	createCustomer(token)

	// 4. 获取客户列表
	fmt.Println("\n4. 获取客户列表...")
	listCustomers(token)

	fmt.Println("\n=== 测试完成 ===")
}

func loginUser() string {
	loginData := LoginRequest{
		Username: "demo",
		Password: "demo",
	}

	reqBody, _ := json.Marshal(loginData)
	resp, err := http.Post("http://localhost:8080/api/v1/user/login.json", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Printf("登录请求失败: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &loginResp)

	return loginResp.Token
}

func createCustomer(token string) {
	customerData := CustomerCreateRequest{
		ClientSessionId: fmt.Sprintf("test-%d", 12345),
		Name:            "测试客户",
		CustomerType:    1,
		Address:         "北京市朝阳区",
		Contacts:        "张三",
		ContactsInfo:    "13800138000",
		Comment:         "重要客户",
		Hidden:          false,
	}

	reqBody, _ := json.Marshal(customerData)
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/customers/add.json", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("创建客户失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("创建客户失败，状态码: %d\n响应: %s\n", resp.StatusCode, string(body))
		return
	}

	var createResp CustomerCreateResponse
	json.Unmarshal(body, &createResp)
	fmt.Printf("创建客户成功！\n")
	fmt.Printf("ID: %d\n", createResp.Id)
	fmt.Printf("名称: %s\n", createResp.Name)
	fmt.Printf("类型: %d\n", createResp.CustomerType)
}

func listCustomers(token string) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/customers/list.json", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("获取客户列表失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("获取客户列表失败，状态码: %d\n响应: %s\n", resp.StatusCode, string(body))
		return
	}

	var customers []CustomerCreateResponse
	json.Unmarshal(body, &customers)

	fmt.Printf("客户列表:\n")
	for i, cust := range customers {
		fmt.Printf("%d. ID:%d 名称:%s 类型:%d 隐藏:%t\n",
			i+1, cust.Id, cust.Name, cust.CustomerType, cust.Hidden)
	}
}