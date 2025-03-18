package com.dev_disitllery.auth.service;

import com.dev_disitllery.auth.model.CustomOAuth2User;
import org.springframework.stereotype.Service;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.security.oauth2.core.OAuth2Error;
import org.springframework.web.client.RestTemplate;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.core.ParameterizedTypeReference;

import java.util.List;
import java.util.Map;

@Service
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        try {
            OAuth2User oAuth2User = super.loadUser(userRequest);

            System.out.println("Attributes de GitHub: " + oAuth2User.getAttributes());

            String name = oAuth2User.getAttribute("login");
            String picture = oAuth2User.getAttribute("avatar_url");

            String email = getEmail(userRequest);

            System.out.println("Name: " + name);
            System.out.println("Email: " + email);
            System.out.println("Picture: " + picture);

            return new CustomOAuth2User(oAuth2User, email);
        } catch (Exception e) {
            OAuth2Error error = new OAuth2Error("oauth2_authentication_error", "Error al cargar el usuario", null);
            throw new OAuth2AuthenticationException(error, e);
        }
    }

    private String getEmail(OAuth2UserRequest userRequest) {
        RestTemplate restTemplate = new RestTemplate();
        HttpHeaders headers = new HttpHeaders();
        headers.setBearerAuth(userRequest.getAccessToken().getTokenValue());
        HttpEntity<String> entity = new HttpEntity<>(headers);

        ResponseEntity<List<Map<String, Object>>> response = restTemplate.exchange(
                "https://api.github.com/user/emails",
                HttpMethod.GET,
                entity,
                new ParameterizedTypeReference<List<Map<String, Object>>>() {}
        );

        if (response.getStatusCode().is2xxSuccessful() && response.getBody() != null) {
            for (Map<String, Object> emailInfo : response.getBody()) {
                if (Boolean.TRUE.equals(emailInfo.get("primary"))) {
                    return (String) emailInfo.get("email");
                }
            }
        }
        return null;
    }
}
