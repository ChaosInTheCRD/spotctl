package internal

import (
   "math/rand"
)
func Generate(n int) string {
    var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
    str := make([]rune, n)
    for i := range str {
        str[i] = chars[rand.Intn(len(chars))]
    }
    return string(str)
}
