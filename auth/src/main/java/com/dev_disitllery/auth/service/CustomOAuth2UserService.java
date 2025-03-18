package com.dev_disitllery.auth.service;

import com.dev_disitllery.auth.model.CustomOAuth2User;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.OAuth2Error;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;

@Service
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        try {
            OAuth2User oAuth2User = super.loadUser(userRequest);

            String email = oAuth2User.getAttribute("email");
            String name = oAuth2User.getAttribute("name");
            String picture = oAuth2User.getAttribute("picture");

            System.out.println("Email: " + email);
            System.out.println("Name: " + name);
            System.out.println("Picture: " + picture);

            return new CustomOAuth2User(oAuth2User);
        } catch (Exception e) {
            OAuth2Error error = new OAuth2Error("oauth2_authentication_error", "Error al cargar el usuario", null);
            throw new OAuth2AuthenticationException(error, e);
        }
    }
}
