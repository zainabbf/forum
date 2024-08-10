## To create the Docker image, we run the following command:
```bash
docker build -t forum-go .
```
## To run the application, we use the following command:
```bash
docker run -p 8080:8080 -it forum-go
```




## Use the enviroment variable to store secrets for authentications . 
.env