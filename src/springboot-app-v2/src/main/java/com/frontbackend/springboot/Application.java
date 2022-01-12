package com.frontbackend.springboot;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

@RestController
@EnableAutoConfiguration
public class Application {

    @RequestMapping("/")
    String index(@RequestHeader HttpHeaders headers) {

        
        headers.forEach((key, value) -> {
            System.out.printf(String.format("Header '%s' = %s\n", key, value));
            // if (key == "x-request-id" ||
            //     key == "x-b3-traceid" ||
            //     key == "x-b3-spanid" ||
            //     key == "x-b3-parentspanid" ||
            //     key == "x-b3-sampled" ||
            //     key == "x-b3-flags" ||
            //     key == "x-ot-span-context") {
            //     requestHeaders.set(key, value);
            // }
        });

        RestTemplate restTemplate = new RestTemplate();
        // String response = restTemplate.getForObject("http://golang-app-data-svc:8080", String.class); 
        // return "response from: springboot app V2, args: " + response;
        ResponseEntity<String> response = restTemplate.exchange(
            "http://golang-app-data-svc:8080", HttpMethod.GET, new HttpEntity<Object>(headers), String.class);

        return "response from: springboot app V2, args: " + response.getBody();
    }

    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
