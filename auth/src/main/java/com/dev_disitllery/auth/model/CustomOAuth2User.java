package com.dev_disitllery.auth.model;

import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.oauth2.core.user.OAuth2User;

import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.Map;

public class CustomOAuth2User implements OAuth2User {

    private final OAuth2User oAuth2User;
    private final String email;
    private final String name;
    private final String picture;
    private final List<Map<String, Object>> publicRepos;

    public CustomOAuth2User(OAuth2User oAuth2User, String email, String name, String picture, List<Map<String, Object>> publicRepos) {
        this.oAuth2User = oAuth2User;
        this.email = email;
        this.name = name;
        this.picture = picture;
        this.publicRepos = publicRepos;
    }

    @Override
    public Map<String, Object> getAttributes() {
        return oAuth2User.getAttributes();
    }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return Collections.emptyList();
    }

    @Override
    public String getName() {
        return this.name != null ? this.name : oAuth2User.getAttribute("login");
    }

    public String getEmail() {
        return email;
    }

    public String getPicture() {
        return picture;
    }

    public Long getGithubId() {
        return oAuth2User.getAttribute("id");
    }

    public String getGithubLogin() {
        return oAuth2User.getAttribute("login");
    }

    public List<Map<String, Object>> getPublicRepos() {
        return publicRepos;
    }

    public int getPublicReposCount() {
        return publicRepos != null ? publicRepos.size() : 0;
    }
}