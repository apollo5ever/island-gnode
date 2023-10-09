package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/deroproject/derohe/cryptography/crypto"
	dero "github.com/deroproject/derohe/rpc"
	"github.com/sirupsen/logrus"
)

var devMode bool

type Bounty struct {
	Name         map[int]string
	Names        []string
	Image        map[int]string
	Images       []string
	Description  map[int]string
	Descriptions []string
	Tagline      map[int]string
	Taglines     []string
	Expiry       time.Time
	Amount       float64
	JudgeList    map[int]string
	Judges       []Name
	JN           float64
	Judge        Name
	JT           float64
	JE           float64
	ExecList     map[int]string
	Execs        []Name
	XN           float64
	X            string
	XT           float64
	XE           float64
	SCID         string
	Index        int
	Status       int
	JF           int
	Initiator    Name
	// Add other attributes as needed
}

type Fundraiser struct {
	Name         map[int]string
	Names        []string
	Image        map[int]string
	Images       []string
	Description  map[int]string
	Descriptions []string
	Tagline      map[int]string
	Taglines     []string
	Goal         float64
	Raised       float64
	Deadline     time.Time
	Claimed      float64
	Index        int
	SCID         string
	Address      string
	Initiator    Name
	Status       int
}

type Tier struct {
	Name         map[int]string
	Names        []string
	Image        map[int]string
	Images       []string
	Description  map[int]string
	Descriptions []string
	Tagline      map[int]string
	Taglines     []string
	Amount       float64
	Interval     float64
	Available    float64
	Address      string
	Index        int
	SCID         string
	Subscribers  []string
}

type Name struct {
	SCID string
	Name string
}

type Island struct {
	SCID        string
	Name        string
	Image       string
	Description string
	Tagline     string
	Bounties    []Bounty
	Fundraisers []Fundraiser
	Tiers       []Tier
	Judging     []Bounty
}

// dReams menu StartGnomon() example

// Name my app
const app_tag = "island_gnode"

// contracts
var registry_scid = "a5daa9a02a81a762c83f3d4ce4592310140586badb4e988431819f47657559f7"

var bounties_scid = "fc2a6923124a07f33c859f201a57159663f087e2f4b163eaa55b0f09bf6de89f"
var fundraisers_scid = "d6ad66e39c99520d4ed42defa4643da2d99f297a506d3ddb6c2aaefbe011f3dc"
var subscriptions_scid = "a4943b10767d3b4b28a0c39fe75303b593b2a8609b07394c803fca1a877716cc"

// Log output
var logger = structures.Logger.WithFields(logrus.Fields{})

var islands []Island

func main() {

	flag.BoolVar(&devMode, "dev", false, "Run in dev mode")
	flag.Parse()
	// create a new Gorilla Mux router
	router := mux.NewRouter()

	// register routes
	//router.HandleFunc("/", getAllIslands).Methods("GET")
	router.HandleFunc("/api/islands", getAllIslands).Methods("GET")
	router.HandleFunc("/api/islands/{id}", getIsland).Methods("GET")

	// Initialize Gnomon fast sync

	// Initialize rpc address to rpc.Daemon var
	if devMode {
		rpc.Daemon.Rpc = "127.0.0.1:20000"
		menu.Gnomes.Fast = false

	} else {
		rpc.Daemon.Rpc = "147.182.177.142:9999"
		menu.Gnomes.Fast = true
	}

	// Initialize logger to Stdout
	menu.InitLogrusLog(runtime.GOOS == "windows")

	rpc.Ping()
	// Check for daemon connection, if daemon is not connected we won't start Gnomon
	if rpc.Daemon.Connect {
		if devMode {
			rpc.SetDaemonClient("http://127.0.0.1:20000")
			rpc.SetWalletClient("http://127.0.0.1:30000", ":")
			bountiesSimBytes, err := ioutil.ReadFile("bounties_sim.txt")
			if err != nil {
				log.Fatal(err)
			}
			bountiesCode := string(bountiesSimBytes)
			time.Sleep(30 * time.Second)
			bounties_scid = InstallContract(bountiesCode, "0")

			registrySimBytes, err := ioutil.ReadFile("registry_sim.txt")
			if err != nil {
				log.Fatal(err)
			}
			registryCode := string(registrySimBytes)
			time.Sleep(30 * time.Second)
			registry_scid = InstallContract(registryCode, "0")

			fundraisersSimBytes, err := ioutil.ReadFile("fundraisers_sim.txt")
			if err != nil {
				log.Fatal(err)
			}
			fundraisersCode := string(fundraisersSimBytes)
			time.Sleep(30 * time.Second)
			fundraisers_scid = InstallContract(fundraisersCode, "0")

			subscriptionsSimBytes, err := ioutil.ReadFile("subscriptions_sim.txt")
			if err != nil {
				log.Fatal(err)
			}
			subscriptionsCode := string(subscriptionsSimBytes)
			time.Sleep(30 * time.Second)
			subscriptions_scid = InstallContract(subscriptionsCode, "0")

			fmt.Println("bounty contract", bounties_scid)
			fmt.Println("registry contract", registry_scid)
			fmt.Println("fundraisers contract", fundraisers_scid)
			fmt.Println("subscriptions contract", subscriptions_scid)
			InstallIsland("apollo", "3")
			InstallIsland("Isle of Wight", "1")
			InstallIsland("Azylem", "2")
		}

		// Initialize NFA search filter and start Gnomon
		filter := []string{"Function Approve(seat Uint64) Uint64", "Function SetTagline(Tagline String) Uint64"}
		menu.StartGnomon(app_tag, "boltdb", filter, 0, 0, nil)

		//var Island = GetIsland("cf530bd98d200171a94bcd6ef1e3ad6348bfa3e6691196e64e93e7953b64a2e4")
		//fmt.Println("main island call", Island)
		//GetAllVars()

		// Exit with ctrl-C
		/* 	var exit bool
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			exit = true
		}() */

		// Gnomon will continue to run if daemon is connected
		/* for !exit && rpc.Daemon.Connect {
			// Example using regexp.Compile
			pattern := "S::PRIVATE-ISLANDS::apollo"
			re, err := regexp.Compile(pattern)
			if err != nil {
				logger.Println(re)
				// Handle the error
			}

			//contracts := menu.Gnomes.GetAllOwnersAndSCIDs()
			//logger.Printf("[%s] Index contains %d contracts\n", app_tag, len(contracts))
			//logger.Println(contracts)
			//logger.Println(menu.Gnomes.GetSCIDValuesByKey("fc2a6923124a07f33c859f201a57159663f087e2f4b163eaa55b0f09bf6de89f", "cf530bd98d200171a94bcd6ef1e3ad6348bfa3e6691196e64e93e7953b64a2e40_E"))
			//_, expiry := menu.Gnomes.GetSCIDValuesByKey("fc2a6923124a07f33c859f201a57159663f087e2f4b163eaa55b0f09bf6de89f", "cf530bd98d200171a94bcd6ef1e3ad6348bfa3e6691196e64e93e7953b64a2e40_E")
			//scid, _ := menu.Gnomes.GetSCIDValuesByKey("a5daa9a02a81a762c83f3d4ce4592310140586badb4e988431819f47657559f7", "S::PRIVATE-ISLANDS::apollo")
			//logger.Println("scid: ", scid)
			//logger.Println("expiry: ", expiry[0])
			//vars := menu.Gnomes.GetAllSCIDVariableDetails("a5daa9a02a81a762c83f3d4ce4592310140586badb4e988431819f47657559f7")
			//logger.Println(vars)
			GetAllVars()
			time.Sleep(3 * time.Second)
			rpc.Ping()
		} */

		// Stop Gnomon
		//menu.Gnomes.Stop(app_tag)
	}
	log.Fatal(http.ListenAndServe(":5000", router))

	//logger.Printf("[%s] Done\n", app_tag)
}

/* func getAllIslands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(islands)
} */

func GetAllVars() []Island {
	//SCID := "a5daa9a02a81a762c83f3d4ce4592310140586badb4e988431819f47657559f7"
	pattern := "S::PRIVATE-ISLANDS::.*"
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Handle the error
	}
	var IslandList []Island

	info := menu.Gnomes.GetAllSCIDVariableDetails(registry_scid)
	if info != nil {
		i := 0
		keys := make([]int, len(info))
		for k := range info {
			keys[i] = int(k)
			i++
		}

		if len(keys) == 0 {
			logger.Warnln("[GetAllVars] No stored heights")
			return []Island{}
		}

		sort.Ints(keys)

		for _, h := range info[int64(keys[len(keys)-1])] {
			//logger.Println("[GetAllVars]", h.Key, h.Value)

			if keyStr, ok := h.Key.(string); ok {

				if re.MatchString(keyStr) {
					logger.Println("Pattern matched for key:", keyStr)
					if valueStr, ok := h.Value.(string); ok {
						logger.Println(valueStr)
						//GetBounty(valueStr, 0)
						IslandList = append(IslandList, GetIsland(valueStr))
					}

				}
			} else {
				logger.Println("Key is not a string.")
			}
		}
		return IslandList
	}
	return []Island{}

}

func GetTier(scid string, index int) Tier {
	info := menu.Gnomes.GetAllSCIDVariableDetails(subscriptions_scid)
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		for k := range info {
			keys[i] = int(k)
			i++
		}
		if len(keys) == 0 {
			fmt.Println("[GetAllVars] No stored heights")
			return Tier{}
		}
		key := scid + strconv.Itoa(index)
		re, err := regexp.Compile(key + ".*")
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return Tier{}
		}
		var Tier Tier = Tier{
			Name:        make(map[int]string),
			Image:       make(map[int]string),
			Description: make(map[int]string),
			Tagline:     make(map[int]string),
			Amount:      float64(0),
			Interval:    float64(0),
			Available:   float64(0),
			Index:       index,
			SCID:        scid,
			Subscribers: make([]string, 0),
		}
		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if re.MatchString(keyStr) {
					parts := strings.Split(keyStr, "_")
					if len(parts) >= 3 {
						versionNumber, _ := strconv.Atoi(parts[len(parts)-1])
						switch parts[len(parts)-2] {
						case "name":
							Tier.Name[versionNumber] = h.Value.(string)
						case "image":
							Tier.Image[versionNumber] = h.Value.(string)
						case "desc":
							Tier.Description[versionNumber] = h.Value.(string)
						case "tagline":
							Tier.Tagline[versionNumber] = h.Value.(string)
						default:
							fmt.Println("TIER TIME DEFAULT PARTY: ??", parts)
							if strings.HasPrefix(parts[0], "dero") {
								fmt.Println("SUBSCRIBER FOUND")
								Tier.Subscribers = append(Tier.Subscribers, parts[0])
							}
						}
					} else {
						switch parts[1] {
						case "Ad":
							address := h.Value.(string)
							Tier.Address = address
						case "Am":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Tier.Amount = amount
						case "Av":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Tier.Available = amount
						case "I":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Tier.Interval = amount

						}
					}
				}
			} else {
				fmt.Println("Key is not a string")
			}
		}
		Tier.Names = MapValuesToSlice(Tier.Name)
		Tier.Images = MapValuesToSlice(Tier.Image)
		Tier.Descriptions = MapValuesToSlice(Tier.Description)
		Tier.Taglines = MapValuesToSlice(Tier.Tagline)

		return Tier

	} else {
		return Tier{}
	}

}

func GetFundraiser(scid string, index int) Fundraiser {
	info := menu.Gnomes.GetAllSCIDVariableDetails(fundraisers_scid)
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		for k := range info {
			keys[i] = int(k)
			i++
		}
		if len(keys) == 0 {
			fmt.Println("[GetFundraiser] No stored heights")
			return Fundraiser{}
		}
		key := scid + strconv.Itoa(index)
		re, err := regexp.Compile(key + ".*")
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return Fundraiser{}
		}
		var Fundraiser Fundraiser = Fundraiser{
			Name:        make(map[int]string),
			Image:       make(map[int]string),
			Description: make(map[int]string),
			Tagline:     make(map[int]string),
			Goal:        float64(0),
			Raised:      float64(0),
			Index:       index,
			SCID:        scid,
			Address:     "",
			Initiator:   getName(scid),
		}
		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if re.MatchString(keyStr) {
					parts := strings.Split(keyStr, "_")
					if len(parts) >= 3 {
						versionNumber, _ := strconv.Atoi(parts[len(parts)-1])
						switch parts[len(parts)-2] {
						case "name":
							Fundraiser.Name[versionNumber] = h.Value.(string)
						case "image":
							Fundraiser.Image[versionNumber] = h.Value.(string)
						case "desc":
							Fundraiser.Description[versionNumber] = h.Value.(string)
						case "tagline":
							Fundraiser.Tagline[versionNumber] = h.Value.(string)
						}
					} else {
						switch parts[1] {
						case "F":
							Fundraiser.Address = h.Value.(string)
						case "D":
							expiryInt, ok := h.Value.(int)
							if !ok {
								expiryFloat := h.Value.(float64)
								expiryInt = int(expiryFloat)
							}
							Fundraiser.Deadline = time.Unix(int64(expiryInt), 0)
						case "G":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Fundraiser.Goal = amount
						case "R":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Fundraiser.Raised = amount
						case "C":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Fundraiser.Claimed = amount

						}
					}
				}
			} else {
				fmt.Println("Key is not a string")
			}
		}
		Fundraiser.Names = MapValuesToSlice(Fundraiser.Name)
		Fundraiser.Images = MapValuesToSlice(Fundraiser.Image)
		Fundraiser.Descriptions = MapValuesToSlice(Fundraiser.Description)
		Fundraiser.Taglines = MapValuesToSlice(Fundraiser.Tagline)

		if Fundraiser.Deadline.After(time.Now().UTC()) {
			Fundraiser.Status = 0
		} else {
			if Fundraiser.Raised < Fundraiser.Goal {
				Fundraiser.Status = 2
			} else {
				Fundraiser.Status = 1
			}
		}

		return Fundraiser

	} else {
		return Fundraiser{}
	}

}

func GetBounty(scid string, index int) Bounty {
	info := menu.Gnomes.GetAllSCIDVariableDetails(bounties_scid)
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		for k := range info {
			keys[i] = int(k)
			i++
		}
		if len(keys) == 0 {
			fmt.Println("[GetAllVars] No stored heights")
			return Bounty{}
		}
		key := scid + strconv.Itoa(index)
		re, err := regexp.Compile(key + ".*")
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return Bounty{}
		}
		var Bounty Bounty = Bounty{
			Name:        make(map[int]string),
			Image:       make(map[int]string),
			Description: make(map[int]string),
			Tagline:     make(map[int]string),
			Amount:      float64(0),
			JudgeList:   make(map[int]string),
			JT:          float64(0),
			ExecList:    make(map[int]string),
			SCID:        scid,
			Index:       index,
			Initiator:   getName(scid),
		}
		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if re.MatchString(keyStr) {
					parts := strings.Split(keyStr, "_")
					if len(parts) >= 3 {
						versionNumber, _ := strconv.Atoi(parts[len(parts)-1])
						switch parts[len(parts)-2] {
						case "name":
							Bounty.Name[versionNumber] = h.Value.(string)
						case "image":
							Bounty.Image[versionNumber] = h.Value.(string)
						case "desc":
							Bounty.Description[versionNumber] = h.Value.(string)
						case "tagline":
							Bounty.Tagline[versionNumber] = h.Value.(string)
						}
					} else {
						switch parts[1] {
						case "E":
							expiryInt, ok := h.Value.(int)
							if !ok {
								expiryFloat := h.Value.(float64)
								expiryInt = int(expiryFloat)
							}
							Bounty.Expiry = time.Unix(int64(expiryInt), 0)
						case "T":
							amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
							if err != nil {
								logger.Println("Failed to convert amount to float64", err)
								continue
							}
							Bounty.Amount = amount
						default:
							if strings.HasPrefix(parts[len(parts)-1], "J") {
								fmt.Println("J DETECTED!!!", parts)
								// Handle JN, JT, J cases
								judgeIndexStr := parts[len(parts)-1]
								fmt.Println("jis", judgeIndexStr)
								judgeIndex, err := strconv.Atoi(judgeIndexStr[1:])
								if err != nil {
									fmt.Println("Invalid JudgeList index:", err)
								} else {
									fmt.Println(judgeIndex)
									fmt.Println(h.Value.(string))
									Bounty.JudgeList[judgeIndex] = h.Value.(string)
								}
								fmt.Println(parts)
								switch parts[len(parts)-1] {
								case "JN":
									JN, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
									if err != nil {
										logger.Println("Failed to convert amount to float64", err)
										continue
									}
									Bounty.JN = JN
								case "JE":
									JE, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
									if err != nil {
										logger.Println("Failed to convert amount to float64", err)
										continue
									}
									Bounty.JE = JE
								case "JT":
									JT, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
									if err != nil {
										logger.Println("Failed to convert amount to float64", err)
										continue
									}
									Bounty.JT = JT
								case "JF":
									Bounty.JF = h.Value.(int)
								case "J":
									// Process J attribute
									Bounty.Judge = getName(h.Value.(string))
								}
							} else if strings.HasPrefix(parts[len(parts)-1], "X") {
								fmt.Println("X DETECTED!!!", parts)
								// Handle XN, XE, XT cases
								xIndexStr := parts[len(parts)-1]
								fmt.Println("xis", xIndexStr)
								xIndex, err := strconv.Atoi(xIndexStr[1:])
								if err != nil {
									fmt.Println("Invalid XList index:", err)
								} else {
									fmt.Println(xIndex)
									fmt.Println(h.Value.(string))
									Bounty.ExecList[xIndex] = h.Value.(string)
								}
								fmt.Println(parts)
								switch parts[len(parts)-1] {
								case "XN":
									XN, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
									if err != nil {
										logger.Println("Failed to convert amount to float64", err)
										continue
									}
									Bounty.XN = XN
								case "XE":
									XE, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
									if err != nil {
										logger.Println("Failed to convert amount to float64", err)
										continue
									}
									Bounty.XE = XE
								case "XT":
									XT, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
									if err != nil {
										logger.Println("Failed to convert amount to float64", err)
										continue
									}
									Bounty.XT = XT
								case "X":
									// Process X attribute
									Bounty.X = h.Value.(string)
								}
							}

						}

					}
				}
			} else {
				fmt.Println("Key is not a string")
			}
		}
		fmt.Println(Bounty.Name)
		fmt.Println(Bounty.JudgeList)
		Bounty.Names = MapValuesToSlice(Bounty.Name)
		Bounty.Images = MapValuesToSlice(Bounty.Image)
		Bounty.Descriptions = MapValuesToSlice(Bounty.Description)
		Bounty.Taglines = MapValuesToSlice(Bounty.Tagline)

		var Judges []string = MapValuesToSlice(Bounty.JudgeList)
		Bounty.Judges = make([]Name, len(Judges))
		for k := range Judges {
			Bounty.Judges[k] = getName(Judges[k])
		}

		var Execs []string = MapValuesToSlice(Bounty.ExecList)
		Bounty.Execs = make([]Name, len(Execs))
		for k := range Execs {
			Bounty.Execs[k] = getName(Execs[k])
		}

		if Bounty.Expiry.Before(time.Now().UTC()) {
			// bounty is expired
			if Bounty.JF == 1 {
				//expired & released
				Bounty.Status = 1
			} else {
				//expired not released
				Bounty.Status = 2
			}
		} else {
			//active
			Bounty.Status = 0
		}

		//Bounty.Judges =
		//Bounty.Execs = MapValuesToSlice(Bounty.ExecList)

		return Bounty

	} else {
		return Bounty{}
	}

}

func GetBounties(scid string) {
	//bountiesSCID := "fc2a6923124a07f33c859f201a57159663f087e2f4b163eaa55b0f09bf6de89f"
	fmt.Println(scid)

	info := menu.Gnomes.GetAllSCIDVariableDetails(bounties_scid)
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		for k := range info {
			keys[i] = int(k)
			i++
		}

		if len(keys) == 0 {
			fmt.Println("[GetAllVars] No stored heights")
			return
		}

		var matchedBounties []Bounty
		for i := 0; ; i++ {
			key := scid + strconv.Itoa(i)

			re, err := regexp.Compile(key + ".*")
			if err != nil {
				fmt.Println("Invalid regex pattern:", err)
				return
			}

			foundMatch := false
			var currentBounty *Bounty = &Bounty{
				Name:        make(map[int]string),
				Image:       make(map[int]string),
				Description: make(map[int]string),
				Tagline:     make(map[int]string),
				Amount:      float64(0),
			}
			for _, h := range info[int64(keys[len(keys)-1])] {
				if keyStr, ok := h.Key.(string); ok {
					if re.MatchString(keyStr) {
						foundMatch = true

						// Split the key string using underscore
						parts := strings.Split(keyStr, "_")
						if len(parts) >= 3 {
							versionNumber, _ := strconv.Atoi(parts[len(parts)-1])
							// Create or get the corresponding bounty object
							/* if currentBounty == nil {
								logger.Println("creating new bounty", scid, i)
								currentBounty = &Bounty{
									Name:        make(map[int]string),
									Image:       make(map[int]string),
									Description: make(map[int]string),
									Tagline:     make(map[int]string),
									Amount:      float64(69),
								}
								logger.Println("Current bounty:", currentBounty)
								matchedBounties = append(matchedBounties, *currentBounty)
								logger.Println("matched bounties", matchedBounties)
							} */
							switch parts[len(parts)-2] {
							case "name":
								currentBounty.Name[versionNumber] = h.Value.(string)
							case "image":
								currentBounty.Image[versionNumber] = h.Value.(string)
							case "desc":
								currentBounty.Description[versionNumber] = h.Value.(string)
							case "tagline":
								currentBounty.Tagline[versionNumber] = h.Value.(string)
								// Add other cases for additional attributes as needed
							}
						} else {
							// Handle other attributes
							switch parts[1] {
							case "E":

								logger.Printf("EXPIRY: %v (type: %T)", h.Value, h.Value)
								/* if currentBounty == nil {
									currentBounty = &Bounty{
										Name:        make(map[int]string),
										Image:       make(map[int]string),
										Description: make(map[int]string),
										Tagline:     make(map[int]string),
										Amount:      float64(0),
									}
									matchedBounties = append(matchedBounties, *currentBounty)
								} */

								// Try to convert to int
								expiryInt, ok := h.Value.(int)
								if !ok {
									// If it's not an int, try to convert to float64
									expiryFloat, ok := h.Value.(float64)
									if !ok {
										logger.Println("Failed to convert expiry to int or float64")
										continue
									}
									// Convert float64 to int64
									expiryInt = int(expiryFloat)
								}

								// Create time.Unix object for Expiry
								currentBounty.Expiry = time.Unix(int64(expiryInt), 0)

							case "T":
								logger.Printf("AMOUNT: %v (type: %T)", h.Value, h.Value)
								amount, err := strconv.ParseFloat(fmt.Sprintf("%v", h.Value), 64)
								if err != nil {
									logger.Println("Failed to convert amount to float64:", err)
									continue
								}
								/* 	if currentBounty == nil {
									currentBounty = &Bounty{
										Name:        make(map[int]string),
										Image:       make(map[int]string),
										Description: make(map[int]string),
										Tagline:     make(map[int]string),
										Amount:      amount,
									}
									matchedBounties = append(matchedBounties, *currentBounty)
								} */
								logger.Println("current bounty amount", scid, i, currentBounty.Amount, amount)
								logger.Println("matched bounties", matchedBounties)
								currentBounty.Amount = amount
								logger.Println("current bounty amount after", scid, i, currentBounty.Amount, amount)
								logger.Println("matched bounties after", matchedBounties)

								// Add other cases for additional attributes as needed
							}

						}
					}
				} else {
					fmt.Println("Key is not a string.")
				}
			}

			if !foundMatch {
				break
			}
			matchedBounties = append(matchedBounties, *currentBounty)
		}

		// Now matchedBounties contains the list of matching bounties
		/* for _, bounty := range matchedBounties {
			fmt.Println("Names:", bounty.Name)
			fmt.Println("Images:", bounty.Image)
			fmt.Println("Taglines:", bounty.Tagline)
			fmt.Println("Expiry:", bounty.Expiry)
			fmt.Println("Amount:", bounty.Amount)
			// Print or use other attributes as needed
		} */
	}
}

func getAllIslands(w http.ResponseWriter, r *http.Request) {
	Islands := GetAllVars()
	IslandsJSON, err := json.Marshal(Islands)
	if err != nil {
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(IslandsJSON)

}

func getIsland(w http.ResponseWriter, r *http.Request) {
	// Extract the "i" value from the URL using gorilla/mux
	vars := mux.Vars(r)
	index := vars["id"]

	islandData := GetIsland(index)

	// Convert the islandData to JSON
	islandJSON, err := json.Marshal(islandData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set CORS headers to allow requests from any origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.Write(islandJSON)
}

func getName(scid string) Name {
	info := menu.Gnomes.GetAllSCIDVariableDetails(registry_scid)
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		for k := range info {
			keys[i] = int(k)
			i++
		}
		if len(keys) == 0 {
			fmt.Println("[getName] No stored heights")
			return Name{}
		}
		var Name Name = Name{}
		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if keyStr == "N::PRIVATE-ISLANDS::"+scid {
					Name.Name = h.Value.(string)
					Name.SCID = scid
					break
				}

			}
		}
		return Name
	}
	return Name{}
}

func RemoveDuplicateBounties(input []Bounty) []Bounty {
	seen := make(map[string]struct{})
	result := []Bounty{}

	for _, bounty := range input {
		if _, exists := seen[bounty.SCID+fmt.Sprint(bounty.Index)]; !exists {
			seen[bounty.SCID+fmt.Sprint(bounty.Index)] = struct{}{}
			result = append(result, bounty)
		}
	}

	return result
}

func getJudging(scid string) []Bounty {
	fmt.Println("boutta get judgin'", scid)
	info := menu.Gnomes.GetAllSCIDVariableDetails(bounties_scid)
	if info != nil {
		i := 0
		keys := make([]int, len(info))
		for k := range info {
			keys[i] = int(k)
			i++
		}
		if len(keys) == 0 {
			logger.Warnln("[getJudging] No stored heights")
			return []Bounty{}
		}
		sort.Ints(keys)
		Bounties := []Bounty{}

		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if valueStr, ok := h.Value.(string); ok {
					if valueStr == scid {
						fmt.Println("found ", keyStr, valueStr)
						parts := strings.Split(keyStr, "_")

						var index, _ = strconv.Atoi(parts[0][64:])
						currentBounty := GetBounty(parts[0][:64], index)
						Bounties = append(Bounties, currentBounty)

					}
				}
			}

		}

		return RemoveDuplicateBounties(Bounties)

	} else {
		return []Bounty{}
	}
}

func GetIsland(scid string) Island {
	info := menu.Gnomes.GetAllSCIDVariableDetails(scid)
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		for k := range info {
			keys[i] = int(k)
			i++
		}
		if len(keys) == 0 {
			fmt.Println("[GetAllVars] No stored heights")
			return Island{}
		}
		var Island Island = Island{}

		Image, err := menu.Gnomes.GetSCIDValuesByKey(scid, "image")
		if err != nil {
			// Handle the error if there is an issue retrieving the values
			log.Fatal(err)
		}
		if len(Image) > 0 {
			Island.Image = Image[0]
		}

		Description, err := menu.Gnomes.GetSCIDValuesByKey(scid, "bio")
		if err != nil {
			// Handle the error if there is an issue retrieving the values
			log.Fatal(err)
		}
		if len(Description) > 0 {
			Island.Description = Description[0]
		}

		Tagline, err := menu.Gnomes.GetSCIDValuesByKey(scid, "tagline")
		if err != nil {
			// Handle the error if there is an issue retrieving the values
			log.Fatal(err)
		}
		if len(Tagline) > 0 {
			Island.Tagline = Tagline[0]
		}
		Island.Name = getName(scid).Name
		var Bounties = []Bounty{}

		for i := 0; ; i++ {
			fmt.Println(i)
			nextBounty := GetBounty(scid, i)

			if nextBounty.Name[0] != "" {
				Bounties = append(Bounties, nextBounty)
				//fmt.Println("bounties now")
				//fmt.Println(Bounties)
			} else {
				break
			}

		}
		var Fundraisers []Fundraiser

		for i := 0; ; i++ {
			fmt.Println(i)
			nextFundraiser := GetFundraiser(scid, i)

			if nextFundraiser.Name[0] != "" {
				Fundraisers = append(Fundraisers, nextFundraiser)
			} else {
				break
			}

		}

		var Tiers []Tier

		for i := 0; ; i++ {
			fmt.Println(i)
			nextTier := GetTier(scid, i)

			if nextTier.Name[0] != "" {
				Tiers = append(Tiers, nextTier)
			} else {
				break
			}

		}

		Island.Judging = getJudging(scid)

		Island.Bounties = Bounties
		Island.Fundraisers = Fundraisers
		Island.Tiers = Tiers
		Island.SCID = scid
		//fmt.Println(Island.Bounties[0].Name)
		//fmt.Println("ISLAND BOY: ", Island)

		return Island

	} else {
		return Island{}
	}

}

func MapValuesToSlice(m map[int]string) []string {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	values := make([]string, 0, len(m))
	for _, k := range keys {
		values = append(values, m[k])
	}

	return values
}

// func getIsland(scid)
//var island
//for i=0 i++
//bounty = getBounty(scid,i)
//island.Bounties append bounty
// same for other two
//return island

//func getBounty(scid,index)

//func getFundraiser(scid, index)

//func getTier(scid, index)

func InstallContract(code string, wallet string) (new_scid string) {

	rpcClientW, ctx, cancel := rpc.SetWalletClient("localhost:3000"+wallet, "")
	defer cancel()

	args := dero.Arguments{}
	txid := dero.Transfer_Result{}

	params := &dero.Transfer_Params{
		Transfers: []dero.Transfer{},
		SC_Code:   code,
		SC_Value:  0,
		SC_RPC:    args,
		Ringsize:  2,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		logger.Errorln("[UploadTokenContract]", err)
		return ""
	}

	logger.Println("InstallContract] Upload TX:", txid)
	rpc.AddLog("Install Contract TX: " + txid.TXID)

	return txid.TXID

}

func InstallIsland(name string, wallet string) (new_scid string) {
	islandBytes, err := ioutil.ReadFile("island.txt")
	if err != nil {
		log.Fatal(err)
	}
	islandCode := string(islandBytes)

	scid := InstallContract(islandCode, wallet)
	time.Sleep(30 * time.Second)
	RegisterIsland(scid, name, wallet)
	fmt.Println(name, scid)

	return scid

}

func RegisterIsland(scid string, name string, wallet string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient("localhost:3000"+wallet, "")
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "RegisterAsset"}
	arg2 := dero.Argument{Name: "scid", DataType: "S", Value: scid}
	arg3 := dero.Argument{Name: "name", DataType: "S", Value: name}
	arg4 := dero.Argument{Name: "collection", DataType: "S", Value: "PRIVATE-ISLANDS"}
	args := dero.Arguments{arg1, arg2, arg3, arg4}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		SCID:   crypto.HashHexToHash(scid),
		Amount: 0,
		Burn:   1,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(registry_scid, "[Register]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     registry_scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		logger.Errorln("[RegisterIsland]", err)
		return
	}

}
