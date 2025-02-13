## Example transaction wrapper

```go
// Create a new user
	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Country:   req.Country,
		Birthday:  birthday,
	}

	// start tx
	result, err := util.WithTransaction(r.Context(), s.DB.GetDB(), func(tx *sql.Tx) (interface{}, error) {

		// Save user to database
		newUser, err := s.UserRepository.CreateUser(r.Context(), tx, &user)
		if err != nil {
			return nil, err
		}

		profile := model.Profile{
			UserID:            newUser.ID,
			DisplayName:       req.DisplayName,
			AvatarVersion:     0,
			Privacy:           req.Privacy,
			FitnessExperience: req.FitnessExperience,
			ExperiencePoints:  0,
		}

		newProfile, err := s.ProfileReposistory.CreateProfile(r.Context(), tx, &profile)
		if err != nil {
			return nil, err
		}

		profileResponseDTO := dto.ProfileResponse{
			ProfileID:         newProfile.ID,
			DisplayName:       newProfile.DisplayName,
			AvatarVersion:     newProfile.AvatarVersion,
			Privacy:           newProfile.Privacy,
			FitnessExperience: newProfile.FitnessExperience,
			ExperiencePoints:  newProfile.ExperiencePoints,
			User: dto.UserResponse{
				UserID:  newUser.ID,
				Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(newUser.FirstName[0])), newUser.LastName),
				Country: newUser.Country,
			},
		}

		//Send Email
		err = s.EmailService.Send("verify_email", "Welcome", []string{newUser.Email}, nil)
		if err != nil {
			util.LogWithContext(
				s.DBLogger,
				slog.LevelError,
				"failed to send email to user",
				map[string]interface{}{
					"userId": newUser.ID,
					"email":  newUser.Email,
				},
				nil)
		}

		return profileResponseDTO, nil
	})

	if err != nil {
		return nil, err
	}

	profileResponseDTO := result.(*dto.ProfileResponse)
	return profileResponseDTO, nil
	//return nil, nil
```