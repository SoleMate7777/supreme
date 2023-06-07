package main

import (
    "fmt"
    "net/http"
    "time"
    "strings"
)

func main() {
    // Get the product URL from the user.
    fmt.Println("Enter the product URL: ")
    var productURL string
    fmt.Scanf("%s", &productURL)

    // Get the product information from the Supreme website.
    response, err := http.Get(productURL)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

    // Check if the product is in stock.
    productInfo := struct {
        InStock bool `json:"in_stock"`
    }{}
    if err := json.NewDecoder(response.Body).Decode(&productInfo); err != nil {
        fmt.Println(err)
        return
    }

    // If the product is in stock, add it to the cart.
    if productInfo.InStock {
        fmt.Println("Product is in stock!")

        // Add the product to the cart.
        addToCartURL := fmt.Sprintf("https://www.supremenewyork.com/cart/add?add=%s", productURL)
        tasks := []struct {
            Name     string
            URL      string
            Proxy    string
            Profile  string
        }{
            {
                Name:     "Task1",
                URL:      addToCartURL,
                Proxy:    "127.0.0.1:8080",
                Profile:  "profile1",
            },
            {
                Name:     "Task2",
                URL:      addToCartURL,
                Proxy:    "127.0.0.1:8081",
                Profile:  "profile2",
            },
            {
                Name:     "Task3",
                URL:      addToCartURL,
                Proxy:    "127.0.0.1:8082",
                Profile:  "profile3",
            },
        }
        for _, task := range tasks {
            // Create a new http client with the proxy and profile.
            client := &http.Client{
                Transport: &http.Transport{
                    Proxy: func(req *http.Request) (*url.URL, error) {
                        return url.Parse(task.Proxy)
                    },
                },
            }

            // Make a request to the Supreme website.
            response, err = client.Get(task.URL)
            if err != nil {
                fmt.Println(err)
                return
            }
            defer response.Body.Close()

            // Check if the request was successful.
            if response.StatusCode == 200 {
                fmt.Println(task.Name, "Request successful!")
            } else {
                fmt.Println(task.Name, "Error making request:", response.StatusCode)
                return
            }

            // Wait for 1 second before checking out.
            time.Sleep(1 * time.Second)

            // Checkout.
            checkoutURL := fmt.Sprintf("https://www.supremenewyork.com/checkout/address")
            response, err = client.Post(checkoutURL, "application/json", nil)
            if err != nil {
                fmt.Println(err)
                return
            }
            defer response.Body.Close()

            // Check if the checkout was successful.
            if response.StatusCode == 200 {
                fmt.Println(task.Name, "Checkout successful!")
            } else {
                fmt.Println(task.Name, "Error checking out:", response.StatusCode)
                return
            }
        }
    } else {
        fmt.Println("Product is out of stock!")
    }
}
