package com.example.todolistapi.service;

import com.example.todolistapi.domain.User;
import com.example.todolistapi.identity.PlatformUser;
import com.example.todolistapi.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class UserService {

    @Autowired
    private final UserRepository userRepository;


    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public void addNewUser(User newUser){

        userRepository.saveAndFlush(newUser);

    }

    public User getUser(String userId){

        List<User> list = userRepository.findUserWithId(userId);

        if(list.isEmpty()){
            return null;
        }
//        Optional fetchedUser = userRepository.f
        return list.get(0);
    }

    public boolean doesUserExistinDB(PlatformUser currentUser){

        User userInDB = getUser(currentUser.getSubject());

        return userInDB != null;
    }
}
