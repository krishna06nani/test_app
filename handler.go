package main

import (
    "ldapconnector"
    "log"
)

func main() {
    serverURI := "ldap://your-ldap-server"
    userDN := "cn=admin,dc=example,dc=com"
    password := "your-password"
    poolSize := 10
    idleTimeout := 300

    pool, err := ldapconnector.NewLDAPPool(serverURI, userDN, password, poolSize)
    if err != nil {
        log.Fatalf("Failed to create LDAP pool: %v", err)
    }

    pool.SetIdleTimeout(idleTimeout)

    // Example usage
    entries, err := pool.Search("dc=example,dc=com", "(cn=John Doe)")
    if err != nil {
        log.Fatalf("LDAP search failed: %v", err)
    }

    for _, entry := range entries {
        log.Printf("DN: %s, CN: %s", entry.DN, entry.GetAttributeValue("cn"))
    }

    pool.CloseAllConnections()
}
