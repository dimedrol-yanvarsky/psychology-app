const lowercase = "abcdefghijklmnopqrstuvwxyz";
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
const numbers = "0123456789";
const specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?";

function shuffleString(string) {
    const array = string.split("");
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array.join("");
}

export function generateStrongPassword() {
    const allChars = lowercase + uppercase + numbers + specialChars;
    const length = Math.floor(Math.random() * 7) + 14;

    let passw = "";
    passw += lowercase[Math.floor(Math.random() * lowercase.length)];
    passw += uppercase[Math.floor(Math.random() * uppercase.length)];
    passw += numbers[Math.floor(Math.random() * numbers.length)];
    passw += specialChars[Math.floor(Math.random() * specialChars.length)];

    for (let i = 4; i < length; i++) {
        passw += allChars[Math.floor(Math.random() * allChars.length)];
    }

    return shuffleString(passw);
}
