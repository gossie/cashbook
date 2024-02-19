package cashbook

import (
	"log"
	"sort"
)

type Person struct {
	Name          string
	sum           float64
	toPay         float64
	OriginalSum   float64
	OriginalToPay float64
}

type Payment struct {
	Id          string
	Amount      float64
	Description string
	Payer       string
}

type Transaction struct {
	From, To *Person
	Amount   float64
}

type Checkout struct {
	CashbookId                  string
	TripName                    string
	Payments                    []*Payment
	TotalCosts, IndividualCosts float64
	Transactions                []*Transaction
	People                      map[string]*Person
}

type Cashbook struct {
	Id       string
	TripName string
	Payments []*Payment
	People   []string
}

func (c *Cashbook) TotalCosts() float64 {
	costs := 0.0
	for _, p := range c.Payments {
		costs += p.Amount
	}
	return costs
}

func (c *Cashbook) IndividualCosts() float64 {
	return c.TotalCosts() / float64(len(c.People))
}

func (c *Cashbook) Checkout() Checkout {
	log.Default().Println("performing checkout of cashbook with id", c.Id)

	peopleMap := make(map[string]*Person)
	for _, name := range c.People {
		peopleMap[name] = &Person{Name: name}
	}

	totalCosts := 0.0
	for _, payment := range c.Payments {
		totalCosts += payment.Amount
		peopleMap[payment.Payer].sum += payment.Amount
	}

	individualCosts := totalCosts / float64(len(peopleMap))

	payList := make([]*Person, 0)
	receiveList := make([]*Person, 0)
	for _, v := range peopleMap {
		if v.sum > individualCosts {
			receiveList = append(receiveList, v)
			v.sum -= individualCosts
			v.OriginalSum = v.sum
		} else {
			v.toPay += individualCosts - v.sum
			v.OriginalToPay = v.toPay
			payList = append(payList, v)
		}
	}

	sort.Slice(receiveList, func(i, j int) bool {
		if receiveList[i].sum == receiveList[j].sum {
			return receiveList[i].Name > receiveList[j].Name
		}
		return receiveList[i].sum > receiveList[j].sum
	})

	sort.Slice(payList, func(i, j int) bool {
		if payList[i].toPay == payList[j].toPay {
			return payList[i].Name > payList[j].Name
		}
		return payList[i].toPay > payList[j].toPay
	})

	transactions := make([]*Transaction, 0)

	if len(payList) > 0 && len(receiveList) > 0 {
		currentPayer := payList[0]
		payList = payList[1:]
		currentReceiver := receiveList[0]
		receiveList = receiveList[1:]

		for currentPayer != nil && currentReceiver != nil {
			if currentPayer.toPay > currentReceiver.sum {
				transactions = append(transactions, &Transaction{currentPayer, currentReceiver, currentReceiver.sum})
				currentPayer.toPay -= currentReceiver.sum
				currentReceiver.sum = 0
				receiveList, currentReceiver = determineNextTurn(receiveList, currentReceiver)
			} else if currentPayer.toPay < currentReceiver.sum {
				transactions = append(transactions, &Transaction{currentPayer, currentReceiver, currentPayer.toPay})
				currentReceiver.sum -= currentPayer.toPay
				currentPayer.toPay = 0
				payList, currentPayer = determineNextTurn(payList, currentPayer)
			} else {
				transactions = append(transactions, &Transaction{currentPayer, currentReceiver, currentReceiver.sum})
				currentPayer.toPay, currentReceiver.sum = 0, 0
				receiveList, currentReceiver = determineNextTurn(receiveList, currentReceiver)
				payList, currentPayer = determineNextTurn(payList, currentPayer)
			}
		}
	}

	return Checkout{
		c.Id,
		c.TripName,
		c.Payments,
		totalCosts,
		individualCosts,
		transactions,
		peopleMap,
	}
}

func determineNextTurn(people []*Person, current *Person) ([]*Person, *Person) {
	if len(people) > 0 {
		return people[1:], people[0]
	}
	return people, nil
}
