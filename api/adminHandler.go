package api


//HTTP Post - /ocms/v2/admin/user/register
func (a *AppContext) registerUser(w http.ResponseWriter, r *http.Request) {
	log.Debug("registerUser() : calling method -")

	var bytes []byte
	var consent helpers.Consent
	err := json.NewDecoder(r.Body).Decode(&consent)
	if err != nil {
		SendError(w, err)
		return
	}
	registerRequest := fabricCAClient.RegistrationRequest{Name: userName, Type: "user", Affiliation: "org1.department1"}
	enrolmentSecret, err := caClient.Register(adminUser, &registerRequest)
	if err != nil {
		t.Fatalf("Error from Register: %s", err)
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", userName, enrolmentSecret)
	// Enrol the previously registered user
	ekey, ecert, err := caClient.Enroll(userName, enrolmentSecret)
	if err != nil {
		t.Fatalf("Error enroling user: %s", err.Error())
	}

	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

//HTTP Post - /ocms/v2/admin/user/enroll
func (a *AppContext) enrollUser(w http.ResponseWriter, r *http.Request) {
	log.Debug("enrollUser() : calling method -")

	var bytes []byte
	var consent helpers.Consent
	err := json.NewDecoder(r.Body).Decode(&consent)
	if err != nil {
		SendError(w, err)
		return
	}
	switch action := consent.Action; action {
	case "create":
		bytes, err = a.createConsent(a.ChainCodeID, consent)
	case "list":
		bytes, err = a.listConsents(a.ChainCodeID, consent.AppID)
	case "get":
		bytes, err = a.getConsent(a.ChainCodeID, consent.AppID, consent.ConsentID)
	case "remove":
		bytes, err = a.unactivateConsent(a.ChainCodeID, consent.AppID, consent.ConsentID)
	case "list4owner":
		bytes, err = a.getConsents4Owner(a.ChainCodeID, consent.AppID, consent.OwnerID)
	case "list4consumer":
		bytes, err = a.getConsents4Consumer(a.ChainCodeID, consent.AppID, consent.ConsumerID)
	case "isconsent":
		bytes, err = a.isConsent(a.ChainCodeID, consent)
	default:
		log.Error("bad action request")
		SendError(w, err)
		return
	}
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

//HTTP Post - /ocms/v2/admin/user/revoke
func (a *AppContext) revokeUser(w http.ResponseWriter, r *http.Request) {
	log.Debug("enrollUser() : calling method -")

	var bytes []byte
	var consent helpers.Consent
	err := json.NewDecoder(r.Body).Decode(&consent)
	if err != nil {
		SendError(w, err)
		return
	}
	switch action := consent.Action; action {
	case "create":
		bytes, err = a.createConsent(a.ChainCodeID, consent)
	case "list":
		bytes, err = a.listConsents(a.ChainCodeID, consent.AppID)
	case "get":
		bytes, err = a.getConsent(a.ChainCodeID, consent.AppID, consent.ConsentID)
	case "remove":
		bytes, err = a.unactivateConsent(a.ChainCodeID, consent.AppID, consent.ConsentID)
	case "list4owner":
		bytes, err = a.getConsents4Owner(a.ChainCodeID, consent.AppID, consent.OwnerID)
	case "list4consumer":
		bytes, err = a.getConsents4Consumer(a.ChainCodeID, consent.AppID, consent.ConsumerID)
	case "isconsent":
		bytes, err = a.isConsent(a.ChainCodeID, consent)
	default:
		log.Error("bad action request")
		SendError(w, err)
		return
	}
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}