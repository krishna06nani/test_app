package main

import (
    "ldapconnector"
    "log"
    "net/http"
    "encoding/json"
)

type SearchRequest struct {
    BaseDN string `json:"base_dn"`
    Filter string `json:"filter"`
}

type SearchResponse struct {
    Entries []map[string]string `json:"entries"`
}

func ldapSearchHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the request body
    var req SearchRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // LDAP connection parameters
    serverURI := "ldap://your-ldap-server"
    userDN := "cn=admin,dc=example,dc=com"
    password := "your-password"
    poolSize := 10
    idleTimeout := 300

    // Create LDAP pool
    pool, err := ldapconnector.NewLDAPPool(serverURI, userDN, password, poolSize)
    if err != nil {
        log.Printf("Failed to create LDAP pool: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    pool.SetIdleTimeout(idleTimeout)

    // Perform LDAP search
    entries, err := pool.Search(req.BaseDN, req.Filter)
    if err != nil {
        log.Printf("LDAP search failed: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Prepare response
    var responseEntries []map[string]string
    for _, entry := range entries {
        entryMap := make(map[string]string)
        entryMap["dn"] = entry.DN
        entryMap["cn"] = entry.GetAttributeValue("cn")
        responseEntries = append(responseEntries, entryMap)
    }

    response := SearchResponse{Entries: responseEntries}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    // Close all connections
    pool.CloseAllConnections()
}

func main() {
    http.HandleFunc("/ldapsearch", ldapSearchHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
