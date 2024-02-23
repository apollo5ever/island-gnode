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

	dero "github.com/deroproject/derohe/rpc"
	"github.com/sirupsen/logrus"
)

//"github.com/deroproject/derohe/cryptography/crypto"

var devMode bool

type Bounty struct {
	Name          string
	Image         string
	Description   string
	Tagline       string
	Expiry        uint64
	Amount        uint64
	Judges        []Judge
	JN            float64
	Judge         Name
	Executer      Name
	JT            float64
	JE            float64
	Execs         []Judge
	XN            float64
	X             string
	XT            float64
	XE            float64
	SCID          string
	Index         uint64
	Status        uint64
	JF            float64
	Initiator     Name
	RecipientList []Recipient
	// Add other attributes as needed
}
type Recipient struct {
	Address string
	Weight  uint64
}

type Judge struct {
	Name  string
	SCID  string
	Index int
}

type Fundraiser struct {
	Name           string
	Image          string
	Description    string
	Tagline        string
	Goal           uint64
	Raised         uint64
	Expiry         uint64
	Claimed        uint64
	Index          int
	SCID           string
	Address        string
	Initiator      Name
	Status         int
	WithdrawlType  string
	ICO            bool
	IcoToken       string
	IcoAmount      uint64
	WithdrawlToken string
}

type Tier struct {
	Name        string
	Image       string
	Description string
	Tagline     string
	Amount      float64
	Interval    float64
	Available   float64
	Address     string
	Index       int
	SCID        string
	Subscribers []string
	Initiator   Name
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
	History     []History
	Bounties    []Bounty
	Fundraisers []Fundraiser
	Tiers       []Tier
	Judging     []Bounty
}

type History struct {
	Height    int
	Attribute string
	Value     string
}

// dReams menu StartGnomon() example

// Name my app
const app_tag = "island_gnode"

// contracts
var registry_scid = "f8a81d0e5c5f9df1f9e41b186f77d1ddbd4daab4e25a380ddde44d66c040da8f"

var bounties_scid = "fc2a6923124a07f33c859f201a57159663f087e2f4b163eaa55b0f09bf6de89f"
var fundraisers_scid = "d6ad66e39c99520d4ed42defa4643da2d99f297a506d3ddb6c2aaefbe011f3dc"
var subscriptions_scid = "ce99dae86c4172378e53be91b4bb2d99f057c1eb24400510621af6002b2b10e3"

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
	router.HandleFunc("/api/contracts", getContracts).Methods("GET")

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

			registrySimBytes, err := ioutil.ReadFile("registry.bas")
			if err != nil {
				log.Fatal(err)
			}
			registryCode := string(registrySimBytes)
			time.Sleep(30 * time.Second)
			registry_scid = InstallContract(registryCode, "0")

			fundraisersSimBytes, err := ioutil.ReadFile("fundraisers.bas")
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
			time.Sleep(30 * time.Second)
			NewCollection(registry_scid)

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
	pattern := "aPRIVATE-ISLANDS.*"
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
	Tier := Tier{
		Index:     index,
		SCID:      scid,
		Initiator: getName(scid),
	}
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
			return Tier
		}
		sort.Ints(keys)

		key := scid + strconv.Itoa(index)
		re, err := regexp.Compile(key + ".*")
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return Tier
		}

		Image, _ := menu.Gnomes.GetSCIDValuesByKey(subscriptions_scid, scid+strconv.Itoa(index)+"Image")

		if len(Image) > 0 {
			Tier.Image = Image[0]
		}

		Description, _ := menu.Gnomes.GetSCIDValuesByKey(subscriptions_scid, scid+strconv.Itoa(index)+"Desc")

		if len(Description) > 0 {
			Tier.Description = Description[0]
		}

		Tagline, _ := menu.Gnomes.GetSCIDValuesByKey(subscriptions_scid, scid+strconv.Itoa(index)+"Tagline")

		if len(Tagline) > 0 {
			Tier.Tagline = Tagline[0]
		}
		Names, _ := menu.Gnomes.GetSCIDValuesByKey(subscriptions_scid, scid+strconv.Itoa(index)+"Name")
		if len(Names) > 0 {
			Tier.Name = Names[0]
		}

		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if re.MatchString(keyStr) {
					parts := strings.Split(keyStr, "_")
					if len(parts) >= 3 {

						switch parts[len(parts)-2] {

						default:
							fmt.Println("TIER TIME DEFAULT PARTY: ??", parts)
							if strings.HasPrefix(parts[0], "dero") {
								fmt.Println("SUBSCRIBER FOUND")
								Tier.Subscribers = append(Tier.Subscribers, parts[0])
							}
						}
					} else if len(parts) == 2 {
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

		return Tier

	} else {
		return Tier
	}

}

func GetFundraiser(scid string, index int) Fundraiser {
	Fundraiser := Fundraiser{
		Index:     index,
		SCID:      scid,
		Initiator: getName(scid),
	}
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
			return Fundraiser
		}
		sort.Ints(keys)
		/* key := scid + strconv.Itoa(index)
		re, err := regexp.Compile(key + ".*")
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return Fundraiser
		} */

		Image, _ := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"Image")

		if len(Image) > 0 {
			Fundraiser.Image = Image[0]
		}

		Description, _ := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"Desc")

		if len(Description) > 0 {
			Fundraiser.Description = Description[0]
		}

		Tagline, _ := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"Tagline")

		if len(Tagline) > 0 {
			Fundraiser.Tagline = Tagline[0]
		}
		Names, _ := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"Name")
		if len(Names) > 0 {
			Fundraiser.Name = Names[0]
		}
		recipient := ""
		recipients, _ := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"_F")
		if len(recipients) > 0 {
			recipient = recipients[0]
		}

		_, Expiry := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"_D")
		if len(Expiry) > 0 {
			Fundraiser.Expiry = Expiry[0]
		}

		_, Raised := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"_R")
		if len(Raised) > 0 {
			Fundraiser.Raised = Raised[0]
		}

		_, Goal := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"_G")
		if len(Goal) > 0 {
			Fundraiser.Goal = Goal[0]
		}

		_, Claimed := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"_C")
		if len(Claimed) > 0 {
			Fundraiser.Claimed = Claimed[0]
		}

		_, WithdrawlType := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"WithdrawlType")
		if len(WithdrawlType) > 0 {
			if WithdrawlType[0] == 0 {
				Fundraiser.WithdrawlType = "address"
				Fundraiser.Address = recipient
			} else {
				Fundraiser.WithdrawlType = "token"
				Fundraiser.WithdrawlToken = recipient
			}

		}

		_, ICO := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"ICO")
		if len(ICO) > 0 {
			if ICO[0] == 0 {
				Fundraiser.ICO = false
			} else {
				Fundraiser.ICO = true
			}

		}

		_, icoAmount := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"icoAmount")
		if len(icoAmount) > 0 {
			Fundraiser.IcoAmount = icoAmount[0]
		}

		icoToken, _ := menu.Gnomes.GetSCIDValuesByKey(fundraisers_scid, scid+strconv.Itoa(index)+"icoToken")
		if len(icoToken) > 0 {
			Fundraiser.IcoToken = icoToken[0]
		}

		/* for _, h := range info[int64(keys[len(keys)-1])] {
		if keyStr, ok := h.Key.(string); ok {
			if re.MatchString(keyStr) {
				parts := strings.Split(keyStr, "_")
				if len(parts) == 2 {

					switch parts[1] {


					case "F":
						Fundraiser.Address = h.Value.(string)
					/* case "D":
						expiryInt, ok := h.Value.(int)
						if !ok {
							expiryFloat := h.Value.(float64)
							expiryInt = int(expiryFloat)
						}
						Fundraiser.Deadline = time.Unix(int64(expiryInt), 0) */
		/* case "G":
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
		} */

		if time.Unix(int64(Fundraiser.Expiry), 0).After(time.Now().UTC()) {
			//active
			Fundraiser.Status = 0
		} else {
			//deadline has past. failure
			if Fundraiser.Raised < Fundraiser.Goal {
				Fundraiser.Status = 2
			} else {
				//success
				Fundraiser.Status = 1
			}
		}

		return Fundraiser

	} else {
		return Fundraiser
	}

}

func GetBounty(scid string, index int) Bounty {
	Bounty := Bounty{
		SCID:          scid,
		Index:         uint64(index),
		Initiator:     getName(scid),
		Judges:        []Judge{},
		Execs:         []Judge{},
		RecipientList: []Recipient{},
	}
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
			return Bounty
		}
		sort.Ints(keys)
		key := scid + strconv.Itoa(index)
		re, err := regexp.Compile(key + ".*")
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return Bounty
		}
		Image, _ := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"Image")

		if len(Image) > 0 {
			Bounty.Image = Image[0]
		}

		Description, _ := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"Desc")

		if len(Description) > 0 {
			Bounty.Description = Description[0]
		}

		Tagline, _ := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"Tagline")

		if len(Tagline) > 0 {
			Bounty.Tagline = Tagline[0]
		}
		Names, _ := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"Name")
		if len(Names) > 0 {
			Bounty.Name = Names[0]
		}
		judge, _ := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"J")
		fmt.Println("JUDGE??", judge)
		if len(judge) > 0 {
			Bounty.Judge = getName(judge[0])
		}

		executer, _ := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"X")

		if len(executer) > 0 {
			Bounty.Executer = getName(executer[0])
		}

		_, Expiry := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"_E")
		if len(Expiry) > 0 {
			Bounty.Expiry = Expiry[0]
		}

		_, Amount := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"_T")
		if len(Amount) > 0 {
			Bounty.Amount = Amount[0]
		}

		for _, h := range info[int64(keys[len(keys)-1])] {
			if keyStr, ok := h.Key.(string); ok {
				if re.MatchString(keyStr) {
					parts := strings.Split(keyStr, "_")
					if len(parts) == 2 {
						switch parts[1] {
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
									Bounty.Judges = append(Bounty.Judges, getJudge(judgeIndex, h.Value.(string)))
								}
								fmt.Println(parts)
								switch parts[len(parts)-1] {
								case "JN":
									JN := h.Value.(float64)

									Bounty.JN = JN
								case "JE":
									JE := h.Value.(float64)

									Bounty.JE = JE
								case "JT":
									JT := h.Value.(float64)

									Bounty.JT = JT
								case "JF":
									Bounty.JF = h.Value.(float64)

								case "J":
									// Process J attribute
									Bounty.Judge = getName(h.Value.(string))
								}
							} else if strings.HasPrefix(parts[len(parts)-1], "X") {

								xIndexStr := parts[len(parts)-1]

								xIndex, err := strconv.Atoi(xIndexStr[1:])
								if err != nil {
									fmt.Println("Invalid XList index:", err)
								} else {
									fmt.Println(xIndex)
									fmt.Println(h.Value.(string))
									Bounty.Execs = append(Bounty.Execs, getJudge(xIndex, h.Value.(string)))

								}
								fmt.Println(parts)
								switch parts[len(parts)-1] {
								case "XN":
									XN := h.Value.(float64)

									Bounty.XN = XN
								case "XE":
									XE := h.Value.(float64)

									Bounty.XE = XE
								case "XT":
									XT := h.Value.(float64)

									Bounty.XT = XT
								case "X":
									// Process X attribute
									Bounty.X = h.Value.(string)
									Bounty.Executer = getName(h.Value.(string))
								}
							} else if strings.HasPrefix(parts[len(parts)-1], "R") {

								rIndexStr := parts[len(parts)-1]

								//rIndex, err := strconv.Atoi(rIndexStr[1:])
								if err != nil {
									fmt.Println("Invalid XList index:", err)
								} else {
									fmt.Println("weight trouble", rIndexStr[1:], scid+strconv.Itoa(index)+"_W"+rIndexStr[1:])
									_, Weight := menu.Gnomes.GetSCIDValuesByKey(bounties_scid, scid+strconv.Itoa(index)+"_W"+rIndexStr[1:])
									fmt.Println("Weight: ", Weight)
									if len(Weight) > 0 {
										fmt.Println("weight >0 see?", Weight)
										Bounty.RecipientList = append(Bounty.RecipientList, Recipient{Address: h.Value.(string), Weight: Weight[0]})
									}
								}

								fmt.Println(parts)
								switch parts[len(parts)-1] {
								case "RN":
									XN := h.Value.(float64)

									Bounty.XN = XN
								case "XE":
									XE := h.Value.(float64)

									Bounty.XE = XE
								case "XT":
									XT := h.Value.(float64)

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

		/* var Judges []string = MapValuesToSlice(Bounty.JudgeList)
		Bounty.Judges = make([]Judge, len(Judges))
		for k := range Judges {
			Bounty.Judges[k].Name = getName(Judges[k])
			Bounty.Judges[k].Index = k
		} */

		if Bounty.JF == 2 {
			//SUCCESS
			Bounty.Status = 1
		} else if time.Unix(int64(Bounty.Expiry), 0).Before(time.Now().UTC()) {
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
		return Bounty
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
			var currentBounty *Bounty = &Bounty{}
			for _, h := range info[int64(keys[len(keys)-1])] {
				if keyStr, ok := h.Key.(string); ok {
					if re.MatchString(keyStr) {
						foundMatch = true

						// Split the key string using underscore
						parts := strings.Split(keyStr, "_")
						if len(parts) >= 3 {

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
								currentBounty.Name = h.Value.(string)
							case "image":
								currentBounty.Image = h.Value.(string)
							case "desc":
								currentBounty.Description = h.Value.(string)
							case "tagline":
								currentBounty.Tagline = h.Value.(string)
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
								expiry := h.Value.(uint64)

								// Create time.Unix object for Expiry
								currentBounty.Expiry = expiry

							case "T":
								logger.Printf("AMOUNT: %v (type: %T)", h.Value, h.Value)
								amount := h.Value.(uint64)

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
	var Name Name = Name{}
	name, err := menu.Gnomes.GetSCIDValuesByKey(registry_scid, "nPRIVATE-ISLANDS"+scid)
	if err != nil {
		// Handle the error if there is an issue retrieving the values
		log.Fatal(err)
	}

	if len(name) > 0 {
		Name.Name = name[0]
		Name.SCID = scid
	}

	return Name
}

func getJudge(index int, scid string) Judge {
	Judge := Judge{}
	name := getName(scid)
	Judge.Name = name.Name
	Judge.SCID = name.SCID
	Judge.Index = index

	return Judge

}

func RemoveDuplicateBounties(input []Bounty) []Bounty {
	fmt.Println("remove duplicate bounties")
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
						fmt.Println("parts", parts)

						var index, _ = strconv.Atoi(parts[0][64:])
						fmt.Println("index", index)
						currentBounty := GetBounty(parts[0][:64], index)
						fmt.Println("currentBounty", currentBounty)
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
	fmt.Println("info", info)
	var Island Island = Island{}
	if info != nil {
		keys := make([]int, len(info))
		i := 0
		fmt.Println("unsorted keys", keys)

		for k := range info {
			keys[i] = int(k)
			i++
			fmt.Println("keys added", i, keys)
		}
		sort.Ints(keys)
		fmt.Println("sorted keys", keys)
		if len(keys) == 0 {
			fmt.Println("[GetAllVars] No stored heights")
			return Island
		}
		for _, k := range keys {
			fmt.Println("k!!", k)
			for _, h := range info[int64(k)] {

				if keyStr, ok := h.Key.(string); ok {
					if keyStr == "C" {
						fmt.Println(keyStr)
					} else {
						hist := History{Attribute: keyStr, Value: h.Value.(string), Height: k}
						fmt.Println("hist", hist)
						Island.History = append(Island.History, hist)
					}

				}
			}
		}

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

			if nextBounty.Name != "" {
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

			if nextFundraiser.Name != "" {
				Fundraisers = append(Fundraisers, nextFundraiser)
			} else {
				break
			}

		}

		var Tiers []Tier

		for i := 0; ; i++ {
			fmt.Println(i)
			nextTier := GetTier(scid, i)

			if nextTier.Name != "" {
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
		return Island
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

func NewCollection(owner string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient("localhost:30000", "")
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "NewCollection"}
	arg2 := dero.Argument{Name: "name", DataType: "S", Value: "PRIVATE-ISLANDS"}
	arg3 := dero.Argument{Name: "owner", DataType: "S", Value: owner}
	arg4 := dero.Argument{Name: "ownerType", DataType: "U", Value: 0}
	arg5 := dero.Argument{Name: "asset1", DataType: "S", Value: ""}
	arg6 := dero.Argument{Name: "asset2", DataType: "S", Value: ""}
	arg7 := dero.Argument{Name: "price1", DataType: "U", Value: 0}
	arg8 := dero.Argument{Name: "price2", DataType: "U", Value: 0}
	arg9 := dero.Argument{Name: "return1", DataType: "U", Value: 0}
	arg10 := dero.Argument{Name: "return2", DataType: "U", Value: 0}
	arg11 := dero.Argument{Name: "collectionType", DataType: "U", Value: 0}
	arg12 := dero.Argument{Name: "twoToken", DataType: "U", Value: 0}
	args := dero.Arguments{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12}
	txid := dero.Transfer_Result{}

	params := &dero.Transfer_Params{
		SC_ID:    registry_scid,
		SC_RPC:   args,
		Ringsize: 2,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		logger.Errorln("[NewCollection]", err)
		return
	}
}

func RegisterIsland(scid string, name string, wallet string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient("localhost:3000"+wallet, "")
	defer cancel()
	//missing args now
	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "RegisterAsset"}
	arg2 := dero.Argument{Name: "scid", DataType: "S", Value: scid}
	arg3 := dero.Argument{Name: "name", DataType: "S", Value: name}
	arg4 := dero.Argument{Name: "collection", DataType: "S", Value: "PRIVATE-ISLANDS"}
	arg5 := dero.Argument{Name: "index", DataType: "U", Value: 0}
	args := dero.Arguments{arg1, arg2, arg3, arg4, arg5}
	txid := dero.Transfer_Result{}

	//t1 := dero.Transfer{
	//	SCID:   crypto.HashHexToHash(scid),
	//	Amount: 0,
	//	Burn:   1,
	//	}

	//t := []dero.Transfer{t1}
	//fee := rpc.GasEstimate(registry_scid, "[Register]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		//Transfers: t,
		SC_ID:    registry_scid,
		SC_RPC:   args,
		Ringsize: 2,
		//Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		logger.Errorln("[RegisterIsland]", err)
		return
	}

}

func getContracts(w http.ResponseWriter, r *http.Request) {
	// Extract the "i" value from the URL using gorilla/mux
	contracts := map[string]string{
		"bounties":      bounties_scid,
		"registry":      registry_scid,
		"subscriptions": subscriptions_scid,
		"fundraisers":   fundraisers_scid,
	}

	// Convert the islandData to JSON
	contractJSON, err := json.Marshal(contracts)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set CORS headers to allow requests from any origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.Write(contractJSON)
}
