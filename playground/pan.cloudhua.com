server {
    listen 80;
    server_name pan.cloudhau.com;

    location / {
        proxy_pass http://192.168.122.11:80;
        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-For $remote_addr;
        client_body_buffer_size 9M;
        client_max_body_size 1000M;
        proxy_buffers 1024 4k;
        proxy_read_timeout 300;
    }
}
