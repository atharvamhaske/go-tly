# Go-tly: A Go-based Tiny URL Shortener.

![Architecture Diagram](./images/GOTLY.png)

Go-tly is a high-performance URL shortener built using Go to levearage go backend and basics of system design, 
featuring:
- gRPC-based key generation service  
- Cloudflare Workers for global-edge redirects  
- Redis queues and cache  
- MongoDB for URL + domain health data  
- ClickHouse for real-time analytics ingestion  
- React/Next.js frontend  
- Domain health checks via Cloudflare API  

