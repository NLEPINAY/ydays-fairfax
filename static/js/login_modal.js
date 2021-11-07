// 🍔 =======  Requêtes AJAX HTTP pour l'envoi des formulaires ======= 🍔

// Routes :
const loginURL = "/login"
const signupURL = "/register"

// Les deux formulaires :
const loginForm = document.getElementById("login-form")
const signupForm = document.getElementById("signup-form")

// Le compteur d'erreur lors du login (compte inexistant) :
let loginCounter = 1

// Les compteurs d'erreur lors du signup (username déjà pris ou email déjà existant) :
let signupUsernameCounter = 1
let signupEmailCounter = 1

// EventListener pour le LOGIN :
loginForm.addEventListener("submit", (e) => {
    e.preventDefault(); // Prevent auto submission and do an AJAX request instead.

    const formData = new FormData(loginForm); // login-identifier, login-password

    const loginData = {
        method: "POST",
        body: formData,
        processData: false, // ! Important
        contentType: false  // ! Important
    }
    fetch(loginURL, loginData)
    .then(res => res.json())
    .then(res => { 
        /** Récupération de la réponse, qui peut prendre 2 formes :
         ** Object { success: true,     message: "Logged in successfully."}
         ** Object { success: false,    message: "No matching account was found."}
         */
        const messageBox = document.getElementById("login-error");

        if(res.success) {
            // ✔️ LOGIN SUCCEEDED :
            messageBox.classList.add("success"); // Classe CSS pour que le message apparaisse en vert
            messageBox.innerHTML = res.message; // "Logged in successfully."
            location.href = "/home"; // Redirection vers la page d'accueil
        } else {
            // ❌ LOGIN FAILED :
            loginCounter < 2 ? messageBox.innerHTML = res.message : messageBox.innerHTML = `${res.message} (${loginCounter})`;
            loginCounter++;
        }
    })
    .catch((error) => console.log(error))
});


// EventListener pour le SIGNUP :
signupForm.addEventListener("submit", (e) => {
    e.preventDefault(); // Prevent auto submission and do an AJAX request instead.

    const formData = new FormData(signupForm); // signup-username, signup-email, signup-password
    
    const signupData = {
        method: "POST",
        body: formData,
        processData: false, // ! Important
        contentType: false  // ! Important
    }
    fetch(signupURL, signupData)
    .then(res => res.json())
    .then(res => { 
        /** Récupération de la réponse, qui peut prendre 3 formes :
         ** Object { success: false,    error: "username",   message: "This username is unavailable."}
         ** Object { success: false,    error: "email",      message: "This email is already registered."}
         ** Object { success: true,     message: "Account successfully created."}
         */

        const messageBoxUsername = document.getElementById("signup-error-username");
        const messageBoxEmail = document.getElementById("signup-error-email");

        messageBoxUsername.innerHTML = "";
        messageBoxEmail.innerHTML = "";

        if(res.success) {
            // ✔️ SIGNUP SUCCEEDED :
            console.log(res.message);
            document.getElementById('modal__container').classList.remove("right__panel__active"); // Ouverture de la partie "Login" pour signifier à l'utilisateur qu'il peut désormais se connecter.
        } else {
            if(res.error == "username") {
                // ❌ SIGNUP FAILED (username unavailable) :
                signupUsernameCounter < 2 ? messageBoxUsername.innerHTML = res.message : messageBoxUsername.innerHTML = `${res.message} (${signupUsernameCounter})`;
                signupUsernameCounter++;
                signupEmailCounter = 1;
            }
            else if (res.error == "email") {
                // ❌ SIGNUP FAILED (email unavailable) :
                signupEmailCounter < 2 ? messageBoxEmail.innerHTML = res.message : messageBoxEmail.innerHTML = `${res.message} (${signupEmailCounter})`;
                signupEmailCounter++;
                signupUsernameCounter = 1;
            }
        }
    })
    .catch((error) => console.log(error))
});

// 🍔 =======  Ouverture / Fermeture de la pop-up ======= 🍔


const main = document.querySelector("main");
const modal = document.querySelector(".modal__wrapper");
const overlay = document.querySelector(".transparent__overlay");

function openModal() {
    main.classList.add("blurry"); // J'ajoute la classe 'blurry' à l'élément main (le fond devient flouté)
    modal.classList.add("active"); // J'ajoute la classe 'active' à la div modal__wraper (la modal apparaît)
    overlay.style.display = "block"; // L'overlay transparent apparaît, donc on pourra cliquer dessus pour fermer la modal
}

function closeModal() {
    main.classList.remove("blurry"); // Je supprime la classe 'blurry' de l'élément main
    modal.classList.remove("active"); // Je supprime la classe 'active' de la div modal__wraper
    overlay.style.display = "none";
}

main.addEventListener("click", () => {
    // Pas sur la landing page, car le bouton Signup se trouve à l'intérieur du tag <main>, donc impossible
    if(location.pathname != "/" && location.pathname != "/login") {
        closeModal();
    }
});

// 🍔 =======  Animation du panneau de la pop-up ======= 🍔

const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('modal__container');

signUpButton.addEventListener('click', () => {
    container.classList.add("right__panel__active");
});

signInButton.addEventListener('click', () => {
    container.classList.remove("right__panel__active");
});

// 🍔 =======  Vérification des données saisies pour s'enregistrer ======= 🍔

var password1 = document.getElementById("signup-password-1");
var password2 = document.getElementById("signup-password-2");
var username = document.getElementById("signup-username");
const button = document.getElementById("signup-button");

const capitals = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowers = "abcdefghijklmnopqrstuvwxyz"
const nums = "0123456789"

// Vérification du match des 2 mots de passe :
function checkPasswords() {
    if (password1.value != password2.value) {
        password2.setCustomValidity("Your passwords do not match.")
        changeToRed(password1)
        changeToRed(password2)
    } else {
        password2.setCustomValidity("")
        changeToNormal(password1)
        changeToNormal(password2)
    }
}

password1.onchange = checkPasswords;
password2.onkeyup = checkPasswords;

// Vérification du pattern du mot de passe :
function hasCapitals(str) {
    for (let i = 0; i < capitals.length; i++) {
        if (str.includes(capitals[i]))
            return true
    }
    return false
}

function hasLowers(str) {
    for (let i = 0; i < lowers.length; i++) {
        if (str.includes(lowers[i]))
            return true
    }
    return false
}

function hasNums(str) {
    for (let i = 0; i < nums.length; i++) {
        if (str.includes(nums[i]))
            return true
    }
    return false
}

function checkPattern() {
    if (!hasCapitals(password1.value) || !hasLowers(password1.value) || !hasNums(password1.value) || password1.value.length < 6 || password1.value.length > 30) {
        password1.setCustomValidity("Your password must contain at least one lowercase character, one uppercase character, one number and be longer than 6 and shorter than 30.")
    } else {
        password1.setCustomValidity("")
    }
}

password1.onchange = checkPattern;
password1.onkeyup = checkPattern;


// Vérification que le username soit valide (uniquement lettres/chiffres/espaces/underscore, et pas 2 espaces consécutifs) :
function checkRegex() {
    for (let i = 0; i < username.value.length - 1; i++) {
        if (!capitals.includes(username.value[i]) && !lowers.includes(username.value[i]) && !nums.includes(username.value[i]) && username.value[i] != " " && username.value[i] != "_" && username.value[i] != ".") {
            return false
        }
    }
    return true
}

function checkConsecutiveSpaces() {
    for (let i = 0; i < username.value.length - 1; i++) {
        if (username.value[i] == " " && username.value[i + 1] == " ") {
            return false
        }
    }
    return true
}

function isValidUsername() {
    if (!checkRegex()) {
        username.setCustomValidity("Your username can only contain letters, numbers, blank spaces, dots and underscores.")
        changeToRed(username, "pattern")
    } else if (!checkConsecutiveSpaces()) {
        username.setCustomValidity("Your username cannot contain several consecutive blank spaces.")
        changeToRed(username)
    } else if (username.value.length < 3 || username.value.length > 20) {
        username.setCustomValidity("Your username must be longer than 3 and shorter than 20 characters.")
        changeToRed(username, "length")
    } else {
        username.setCustomValidity("")
        changeToNormal(username)
    }
}

function changeToRed(element, str) {
    element.classList.add("invalid__input");
    if (element == username) {
        // Désactivation du bouton 'Continue'
        disables(button)
        if (str == "length") {
            document.getElementsByClassName("invalid__username")[0].style.display = "block"
            document.getElementsByClassName("invalid__username")[1].style.display = "none"
        }
        if (str == "pattern") {
            document.getElementsByClassName("invalid__username")[1].style.display = "block"
            document.getElementsByClassName("invalid__username")[0].style.display = "none"
        }
    }
    if (element == password1) {
        document.querySelector(".invalid__password").style.display = "block"
        disables(button)
    }
}

function changeToNormal(element) {
    enables(button)
    element.classList.remove("invalid__input");
    if (element == username) {
        document.getElementsByClassName("invalid__username")[0].style.display = "none"
        document.getElementsByClassName("invalid__username")[1].style.display = "none"
    }
    if (element == password1) {
        document.querySelector(".invalid__password").style.display = "none"
    }
}

function disables(element) {
    // Désactivation du bouton "Continue" :
    element.disabled = true
    element.style.opacity = 0.4
    element.style.pointerEvents = "none"
}

function enables(element) {
    // Activation du bouton "Continue" :
    element.disabled = false
    element.style.opacity = 1
    element.style.pointerEvents = "initial"
}

username.onchange = isValidUsername;
username.onkeyup = isValidUsername;

// Activation des div "requirement" lorsque le mot de passe est valide :
let letterTyped = "";
const requirements = document.getElementsByClassName("requirement")
const ticks = document.querySelectorAll("svg")

password1.addEventListener("input", (e) => {
    letterTyped = e.target.value;
    if (!isLongEnough()) {
        // Désactivation de la 1ère div :
        requirements[0].style.opacity = 0.35
        ticks[0].style.opacity = 0
    } else {
        // Activation de la 1ère div :
        requirements[0].style.opacity = 1
        ticks[0].style.opacity = 1
    }
    if (!hasCorrectPattern()) {
        // Désactivation de la 2nde div :
        requirements[1].style.opacity = 0.35
        ticks[1].style.opacity = 0
    } else {
        // Activation de la 2nde div :
        requirements[1].style.opacity = 1
        ticks[1].style.opacity = 1
    }
    if (!isLongEnough() || !hasCorrectPattern()) {
        // Désactivation du bouton "Continue" :
        disables(button)
    } else {
        // Activation du bouton "Continue" :
        enables(button)
    }
});

function hasCorrectPattern() {
    if (!hasCapitals(letterTyped) || !hasLowers(letterTyped) || !hasNums(letterTyped)) {
        return false
    } else {
        return true
    }
}

function isLongEnough() {
    if (letterTyped.length < 6 || letterTyped.length > 30) {
        return false
    } else {
        return true
    }
}