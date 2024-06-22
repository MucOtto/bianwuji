# 边无际笔试

## 1、编写golang程序

```go
package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getMeasurements(url string) ([]float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	measurements := make([]float64, 0)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		for _, v := range values {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			measurements = append(measurements, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return measurements, nil
}

func main() {
	serviceURL := os.Getenv("SERVICE_URL")
	if serviceURL == "" {
		log.Fatalf("SERVICE_URL environment variable is required")
	}

	intervalStr := os.Getenv("INTERVAL_SECONDS")
	if intervalStr == "" {
		log.Fatalf("INTERVAL_SECONDS environment variable is required")
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			measurements, err := getMeasurements(serviceURL)
			if err != nil {
				log.Printf("Error getting measurements: %v\n", err)
				continue
			}

			if len(measurements) == 0 {
				log.Println("No valid measurements retrieved")
				continue
			}

			sum := 0.0
			for _, value := range measurements {
				sum += value
			}
			average := sum / float64(len(measurements))
			log.Printf("Average measurement: %.2f\n", average)
		}
	}
}

```

## 2、编写dockerfile

```dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app bishi.go

CMD ["./app", "http://deviceshifu-plate-reader.deviceshifu.svc.cluster.local/get_measurement"]

```



## 3、编写k8s  yaml

轮训时间可以在这里进行修改

### deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: enzyme-reader
spec:
  replicas: 1
  selector:
    matchLabels:
      app: enzyme-reader
  template:
    metadata:
      labels:
        app: enzyme-reader
    spec:
      containers:
        - name: enzyme-reader
          image: xiaoshi1980/enzyme-reader:latest
          env:
            - name: SERVICE_URL
              value: "http://deviceshifu-plate-reader.deviceshifu.svc.cluster.local/get_measurement"
            - name: INTERVAL_SECONDS
              value: "10"
          ports:
            - containerPort: 8080


```



### service:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: enzyme-reader-service
spec:
  selector:
    app: enzyme-reader
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080

```



## 4、构建docker镜像

```bash
 docker build -t xiaoshi1980/enzyme-reader:latest .
```

## 5、部署服务到k8s集群

```bash
kubectl apply -f deployment.yaml 
```



## 6、创建service

```bash
kubectl apply -f service.yaml
```



## 7、查看服务

<img width="669" alt="image" src="https://github.com/MucOtto/bianwuji/assets/122969909/cc171827-285f-4c4f-baad-bf91f58ef6c6">

