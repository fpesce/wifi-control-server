package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	routerURLEnv = "ROUTER_URL"
	usernameEnv  = "USERNAME"
	passwordEnv  = "PASSWORD"
)

func main() {
	routerURL := os.Getenv(routerURLEnv)
	username := os.Getenv(usernameEnv)
	password := os.Getenv(passwordEnv)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, r, routerURL, username, password)
	})

	http.HandleFunc("/wifi-control", func(w http.ResponseWriter, r *http.Request) {
		wifiControlHandler(w, r, routerURL, username, password)
	})


	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wifiControlHandler(w http.ResponseWriter, r *http.Request, routerURL, username, password string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Command not provided", http.StatusBadRequest)
		return
	}

	var postData string
	if command == "on" {
		postData = "submit_flag=wlan_adv&wl_rts=2347&wl_frag=2346&wl_enable_shortpreamble=0&wl_tx_ctrl=100&wl_enable_router=1&endis_pin=0&hid_protect_enable=1&hid_super_wifi=0&hid_super_wifi_an=0&wla_rts=2347&wla_frag=&wla_enable_shortpreamble=0&wla_tx_ctrl=100&wla_enable_router=0&wds_change_ip=still_lanip&enable_router=ignore&enable_ssid_broadcast=ignore&endis_wsc_config=5&endis_wsc_config_a=5&wladv_endis_wmm=1&wladv_endis_wmm_a=1&hid_dyn_get_ip=&hid_ap_ipaddr=&hid_ap_subnet=&hid_ap_gateway=&ap_dnsaddr1=&ap_dnsaddr2=&wladv_enable_schedule=0&wladv_enable_schedule_a=0&wladv_schedule_type=&wladv_schedule_edit_num=&wladv_schedule_delete_num=&hid_enable_coexist=0&hid_wla_beamforming=1&hid_wla_mu_mimo=1&hid_enable_atf=0&hid_wla_ht160=0&enable_ap=1&enable_coexistence=on&wmm_enable=1&frag=2346&rts=2347&enable_shortpreamble=automatic&tx_power_ctrl=100&wmm_enable_a=1&frag_an=2346&rts_an=2347&enable_shortpreamble_an=automatic&tx_power_ctrl_an=100&enable_implicit_beamforming=0&enable_mu="
	} else if command == "off" {
		postData = "submit_flag=wlan_adv&wl_rts=2347&wl_frag=2346&wl_enable_shortpreamble=0&wl_tx_ctrl=100&wl_enable_router=0&endis_pin=0&hid_protect_enable=1&hid_super_wifi=0&hid_super_wifi_an=0&wla_rts=2347&wla_frag=&wla_enable_shortpreamble=0&wla_tx_ctrl=100&wla_enable_router=1&wds_change_ip=still_lanip&enable_router=ignore&enable_ssid_broadcast=ignore&endis_wsc_config=5&endis_wsc_config_a=5&wladv_endis_wmm=1&wladv_endis_wmm_a=1&hid_dyn_get_ip=&hid_ap_ipaddr=&hid_ap_subnet=&hid_ap_gateway=&ap_dnsaddr1=&ap_dnsaddr2=&wladv_enable_schedule=0&wladv_enable_schedule_a=0&wladv_schedule_type=&wladv_schedule_edit_num=&wladv_schedule_delete_num=&hid_enable_coexist=0&hid_wla_beamforming=1&hid_wla_mu_mimo=1&hid_enable_atf=0&hid_wla_ht160=0&enable_ap=1&enable_coexistence=on&wmm_enable=1&frag=2346&rts=2347&enable_shortpreamble=automatic&tx_power_ctrl=100&wmm_enable_a=1&frag_an=2346&rts_an=2347&enable_shortpreamble_an=automatic&tx_power_ctrl_an=100&enable_implicit_beamforming=0&enable_mu="
	} else {
		http.Error(w, "Invalid command", http.StatusBadRequest)
		return
	}

	timestamp, err := getTimestamp(routerURL + "/WLG_adv.htm", username, password)
	if err != nil {
		http.Error(w, "Failed to get timestamp", http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	err = sendPostRequest(routerURL+"/apply.cgi?/WLG_adv.htm%20timestamp="+timestamp, postData, username, password)
	if err != nil {
		http.Error(w, "Failed to send request to router", http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Fprint(w, "Success")
}

func getTimestamp(urlStr, username, password string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return "", err
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`var ts="(\d+)";`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		return "", fmt.Errorf("Timestamp not found")
	}

	return matches[1], nil
}

func sendPostRequest(urlStr, data, username, password string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func serveHTML(w http.ResponseWriter, r *http.Request, routerURL, username, password string) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	wifiStatus, err := getWiFiStatus(routerURL + "/basic_home.htm", username, password)
	if err != nil {
		http.Error(w, "Failed to get WiFi status", http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}


	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>WiFi Control</title>
	<script>
		function sendCommand(command) {
			const xhr = new XMLHttpRequest();
			xhr.open("POST", "/wifi-control");
			xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
			xhr.onreadystatechange = function() {
				if (xhr.readyState === XMLHttpRequest.DONE) {
					alert("Result: " + xhr.responseText);
				}
			};
			xhr.send("command=" + encodeURIComponent(command));
		}
	</script>
</head>
<body>
	<h1>WiFi Control</h1>
	<p>WiFi Status: <strong>%s</strong></p>
	<button onclick="sendCommand('on')">Turn WiFi On</button>
	<button onclick="sendCommand('off')">Turn WiFi Off</button>
</body>
</html>
`, wifiStatus)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func getWiFiStatus(urlStr, username, password string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return "", err
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`var enabled_wps="(\d)";`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		return "", fmt.Errorf("WiFi status not found")
	}

	wifiStatus := "Down"
	if matches[1] == "1" {
		wifiStatus = "Up"
	}

	return wifiStatus, nil
}

