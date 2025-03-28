package com.dev_disitllery.auth.service;

import com.dev_disitllery.auth.model.CustomOAuth2User;
import com.dev_disitllery.auth.model.User;
import com.dev_disitllery.auth.repository.UserRepository;
import org.springframework.stereotype.Service;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
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

    private final JwtService jwtService;
    private final UserRepository userRepository;

    public CustomOAuth2UserService(JwtService jwtService, UserRepository userRepository) {
        this.jwtService = jwtService;
        this.userRepository = userRepository;
    }

    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2User oAuth2User = super.loadUser(userRequest);

        Map<String, Object> attributes = oAuth2User.getAttributes();
        String name = (String) attributes.get("name");
        String login = (String) attributes.get("login");
        String avatarUrl = (String) attributes.get("avatar_url");

        String email = getEmail(userRequest);

        User user = userRepository.findByEmail(email)
                .orElse(new User());

        user.setEmail(email);
        user.setName(name != null ? name : login);
        userRepository.save(user);

        String token = jwtService.generateToken(email);
        System.out.println("Token generado para " + email + ": " + token);

        return new CustomOAuth2User(oAuth2User, email, name, avatarUrl);
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
        throw new OAuth2AuthenticationException("No se pudo obtener el email del usuario");
    }
}