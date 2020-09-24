/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package storage

import (
	"context"
	"errors"
	"health/core"
	"health/core/model"
	"io/ioutil"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	configs        *collectionWrapper
	users          *collectionWrapper
	providers      *collectionWrapper
	locations      *collectionWrapper
	ctests         *collectionWrapper
	emanualtests   *collectionWrapper
	resources      *collectionWrapper
	faq            *collectionWrapper
	news           *collectionWrapper
	estatus        *collectionWrapper
	ehistory       *collectionWrapper
	counties       *collectionWrapper
	testtypes      *collectionWrapper
	rules          *collectionWrapper
	symptomgroups  *collectionWrapper //old
	symptomrules   *collectionWrapper //old
	symptoms       *collectionWrapper
	crules         *collectionWrapper
	traceexposures *collectionWrapper
	accessrules    *collectionWrapper
	uinoverrides   *collectionWrapper

	listener core.StorageListener
}

func (m *database) start() error {
	log.Println("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)
	configs := &collectionWrapper{database: m, coll: db.Collection("configs")}
	err = m.applyConfigsChecks(configs)
	if err != nil {
		return err
	}
	users := &collectionWrapper{database: m, coll: db.Collection("users")}
	err = m.applyUsersChecks(users)
	if err != nil {
		return err
	}
	providers := &collectionWrapper{database: m, coll: db.Collection("providers")}
	err = m.applyProvidersChecks(providers)
	if err != nil {
		return err
	}
	locations := &collectionWrapper{database: m, coll: db.Collection("locations")}
	err = m.applyLocationsChecks(locations)
	if err != nil {
		return err
	}
	ctests := &collectionWrapper{database: m, coll: db.Collection("ctests")}
	err = m.applyCTestsChecks(ctests)
	if err != nil {
		return err
	}
	manualtests := &collectionWrapper{database: m, coll: db.Collection("manualtests")}
	err = m.applyManualTestsChecks(manualtests)
	if err != nil {
		return err
	}
	emanualtests := &collectionWrapper{database: m, coll: db.Collection("emanualtests")}
	err = m.applyEManualTestsChecks(emanualtests)
	if err != nil {
		return err
	}
	resources := &collectionWrapper{database: m, coll: db.Collection("resources")}
	err = m.applyResourcesChecks(resources)
	if err != nil {
		return err
	}
	faq := &collectionWrapper{database: m, coll: db.Collection("faq")}
	err = m.applyFAQChecks(faq)
	if err != nil {
		return err
	}
	news := &collectionWrapper{database: m, coll: db.Collection("news")}
	err = m.applyNewsChecks(news)
	if err != nil {
		return err
	}
	status := &collectionWrapper{database: m, coll: db.Collection("status")}
	err = m.applyStatusChecks(status)
	if err != nil {
		return err
	}
	estatus := &collectionWrapper{database: m, coll: db.Collection("estatus")}
	err = m.applyEStatusChecks(estatus)
	if err != nil {
		return err
	}
	history := &collectionWrapper{database: m, coll: db.Collection("history")}
	err = m.applyHistoryChecks(history)
	if err != nil {
		return err
	}
	ehistory := &collectionWrapper{database: m, coll: db.Collection("ehistory")}
	err = m.applyEHistoryChecks(ehistory)
	if err != nil {
		return err
	}
	counties := &collectionWrapper{database: m, coll: db.Collection("counties")}
	err = m.applyCountiesChecks(counties)
	if err != nil {
		return err
	}
	testtypes := &collectionWrapper{database: m, coll: db.Collection("testtypes")}
	err = m.applyTestTypesChecks(testtypes)
	if err != nil {
		return err
	}
	rules := &collectionWrapper{database: m, coll: db.Collection("rules")}
	err = m.applyRulesChecks(rules)
	if err != nil {
		return err
	}
	symptomgroups := &collectionWrapper{database: m, coll: db.Collection("symptomgroups")}
	err = m.applySymptomGroupsChecks(symptomgroups)
	if err != nil {
		return err
	}
	symptomrules := &collectionWrapper{database: m, coll: db.Collection("symptomrules")}
	err = m.applySymptomRulesChecks(symptomrules)
	if err != nil {
		return err
	}
	symptoms := &collectionWrapper{database: m, coll: db.Collection("symptoms")}
	err = m.applySymptomsChecks(symptoms)
	if err != nil {
		return err
	}
	crules := &collectionWrapper{database: m, coll: db.Collection("crules")}
	err = m.applyCRulesChecks(crules, counties)
	if err != nil {
		return err
	}
	traceexposures := &collectionWrapper{database: m, coll: db.Collection("traceexposures")}
	err = m.applyTraceExposuresChecks(traceexposures)
	if err != nil {
		return err
	}
	accessrules := &collectionWrapper{database: m, coll: db.Collection("accessrules")}
	err = m.applyAccessRulesChecks(accessrules)
	if err != nil {
		return err
	}
	uinoverrides := &collectionWrapper{database: m, coll: db.Collection("uinoverrides")}
	err = m.applyUINOverridesChecks(uinoverrides)
	if err != nil {
		return err
	}

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	m.configs = configs
	m.users = users
	m.providers = providers
	m.locations = locations
	m.ctests = ctests
	m.emanualtests = emanualtests
	m.resources = resources
	m.faq = faq
	m.news = news
	m.estatus = estatus
	m.ehistory = ehistory
	m.counties = counties
	m.testtypes = testtypes
	m.rules = rules
	m.symptomgroups = symptomgroups
	m.symptomrules = symptomrules
	m.symptoms = symptoms
	m.crules = crules
	m.traceexposures = traceexposures
	m.accessrules = accessrules
	m.uinoverrides = uinoverrides

	//watch for config changes
	go m.configs.Watch(nil)

	return nil
}

func (m *database) applyConfigsChecks(configs *collectionWrapper) error {
	log.Println("apply configs checks.....")

	log.Println("consfigs checks passed")
	return nil
}

func (m *database) applyUsersChecks(users *collectionWrapper) error {
	log.Println("apply users checks.....")

	//add external id index - unique
	err := users.AddIndex(bson.D{primitive.E{Key: "external_id", Value: 1}}, true)
	if err != nil {
		return err
	}

	//add shibboleth index
	err = users.AddIndex(bson.D{primitive.E{Key: "shibboleth_auth.uiucedu_uin", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add uuid index
	err = users.AddIndex(bson.D{primitive.E{Key: "uuid", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add re_post index
	err = users.AddIndex(bson.D{primitive.E{Key: "re_post", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("users checks passed")
	return nil
}

func (m *database) applyProvidersChecks(providers *collectionWrapper) error {
	log.Println("apply providers checks.....")

	log.Println("providers checks passed")
	return nil
}

func (m *database) applyLocationsChecks(locations *collectionWrapper) error {
	log.Println("apply locations checks.....")

	//add indexes
	err := locations.AddIndex(bson.D{primitive.E{Key: "provider_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = locations.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("locations checks passed")
	return nil
}

func (m *database) applyCTestsChecks(ctests *collectionWrapper) error {
	log.Println("apply ctests checks.....")

	//add indexes
	err := ctests.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = ctests.AddIndex(bson.D{primitive.E{Key: "provider_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = ctests.AddIndex(bson.D{primitive.E{Key: "order_number", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("ctests checks passed")
	return nil
}

func (m *database) applyManualTestsChecks(manualtests *collectionWrapper) error {
	log.Println("apply manualtests checks.....")

	//add indexes
	err := manualtests.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = manualtests.AddIndex(bson.D{primitive.E{Key: "location_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = manualtests.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = manualtests.AddIndex(bson.D{primitive.E{Key: "status", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("manualtests checks passed")
	return nil
}

func (m *database) applyEManualTestsChecks(emanualtests *collectionWrapper) error {
	log.Println("apply emanualtests checks.....")

	//add indexes
	err := emanualtests.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = emanualtests.AddIndex(bson.D{primitive.E{Key: "location_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = emanualtests.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = emanualtests.AddIndex(bson.D{primitive.E{Key: "status", Value: 1}}, false)
	if err != nil {
		return err
	}

	// Remove all verified manual tests as we already do not keep them
	// First check their count
	verifiedFilter := bson.D{primitive.E{Key: "status", Value: "verified"}}
	var items []*eManualTest
	err = emanualtests.Find(verifiedFilter, &items, nil)
	if err != nil {
		return err
	}
	if items != nil && len(items) > 0 {
		log.Printf("there there are %d verified items, so remove them\n", len(items))

		result, err := emanualtests.DeleteMany(verifiedFilter, nil)
		if err != nil {
			return err
		}
		if result == nil {
			return errors.New("delete result is nil for some reasons")
		}
		log.Printf("%d items were removed\n", result.DeletedCount)
	} else {
		log.Println("there is no verified manual test items, so do nothing")
	}

	log.Println("emanualtests checks passed")
	return nil
}

func (m *database) applyResourcesChecks(resources *collectionWrapper) error {
	log.Println("apply resources checks.....")

	log.Println("resources checks passed")
	return nil
}

func (m *database) applyFAQChecks(faq *collectionWrapper) error {
	log.Println("apply faq checks.....")

	log.Println("faq checks passed")
	return nil
}

func (m *database) applyNewsChecks(news *collectionWrapper) error {
	log.Println("apply news checks.....")

	//add index
	err := news.AddIndex(bson.D{primitive.E{Key: "date", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("news checks passed")
	return nil
}

func (m *database) applyStatusChecks(status *collectionWrapper) error {
	log.Println("apply status checks.....")

	//add index - unique
	err := status.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, true)
	if err != nil {
		return err
	}

	log.Println("status checks passed")
	return nil
}

func (m *database) applyEStatusChecks(estatus *collectionWrapper) error {
	log.Println("apply estatus checks.....")

	//add user_id index
	err := estatus.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add app_version index
	err = estatus.AddIndex(bson.D{primitive.E{Key: "app_version", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("estatus checks passed")
	return nil
}

func (m *database) applyHistoryChecks(history *collectionWrapper) error {
	log.Println("apply history checks.....")

	//add index
	err := history.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add index
	err = history.AddIndex(bson.D{primitive.E{Key: "date", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("history checks passed")
	return nil
}

func (m *database) applyEHistoryChecks(ehistory *collectionWrapper) error {
	log.Println("apply ehistory checks.....")

	//add index
	err := ehistory.AddIndex(bson.D{primitive.E{Key: "user_id", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add index
	err = ehistory.AddIndex(bson.D{primitive.E{Key: "date", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("ehistory checks passed")
	return nil
}

func (m *database) applyCountiesChecks(counties *collectionWrapper) error {
	log.Println("apply counties checks.....")

	//add index
	err := counties.AddIndex(bson.D{primitive.E{Key: "guidelines.id", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add index
	err = counties.AddIndex(bson.D{primitive.E{Key: "county_statuses.id", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("counties checks passed")
	return nil
}

func (m *database) applyTestTypesChecks(testTypes *collectionWrapper) error {
	log.Println("apply testTypes checks.....")

	//add index
	err := testTypes.AddIndex(bson.D{primitive.E{Key: "results._id", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("testTypes checks passed")
	return nil
}

func (m *database) applyRulesChecks(rules *collectionWrapper) error {
	log.Println("apply rules checks.....")

	//add index
	err := rules.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = rules.AddIndex(bson.D{primitive.E{Key: "test_type_id", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("rules checks passed")
	return nil
}

func (m *database) applySymptomGroupsChecks(symptomGroups *collectionWrapper) error {
	log.Println("apply symptomGroups checks.....")

	//1. add index
	err := symptomGroups.AddIndex(bson.D{primitive.E{Key: "symptoms.id", Value: 1}}, false)
	if err != nil {
		return err
	}

	//2. check if need to add the two groups
	filter := bson.D{}
	var result []symptomGroup
	err = symptomGroups.Find(filter, &result, nil)
	if err != nil {
		return err
	}
	hasData := result != nil && len(result) > 0
	if !hasData {
		log.Println("there is no symptoms groups data, so create a default one")

		gr1ID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		gr1 := symptomGroup{ID: gr1ID.String(), Name: "gr1"}
		_, err = symptomGroups.InsertOne(&gr1)
		if err != nil {
			return err
		}

		gr2ID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		gr2 := symptomGroup{ID: gr2ID.String(), Name: "gr2"}
		_, err = symptomGroups.InsertOne(&gr2)
		if err != nil {
			return err
		}
	} else {
		log.Println("there is symptoms groups data, so do nothing")
	}

	log.Println("symptomGroups checks passed")
	return nil
}

func (m *database) applySymptomRulesChecks(symptomRules *collectionWrapper) error {
	log.Println("apply symptomRules checks.....")

	//add index
	err := symptomRules.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, true)
	if err != nil {
		return err
	}

	log.Println("symptomRules checks passed")
	return nil
}

func (m *database) applySymptomsChecks(symptoms *collectionWrapper) error {
	log.Println("apply symptoms checks.....")

	//add index
	err := symptoms.AddIndex(bson.D{primitive.E{Key: "app_version", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add initial data for version 2.6 if not added
	filter := bson.D{primitive.E{Key: "app_version", Value: "2.6"}}
	var items []*model.Symptom
	err = symptoms.Find(filter, &items, nil)
	if err != nil {
		return err
	}
	if len(items) <= 0 {
		log.Println("there are no symptoms for version 2.6, so we need to add initial data")

		data, err := ioutil.ReadFile("./driven/storage/symptoms_2.6.json")
		if err != nil {
			return err
		}
		d := model.Symptoms{AppVersion: "2.6", Items: string(data)}
		_, err = symptoms.InsertOne(&d)
		if err != nil {
			return err
		}
	} else {
		log.Println("there are symptoms for version 2.6, so nothing to do")
	}

	log.Println("symptoms checks passed")
	return nil
}

func (m *database) applyCRulesChecks(cRules *collectionWrapper, counties *collectionWrapper) error {
	log.Println("apply CRules checks.....")

	//add indexes
	err := cRules.AddIndex(bson.D{primitive.E{Key: "app_version", Value: 1}}, false)
	if err != nil {
		return err
	}

	err = cRules.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add initial data for version 2.6 and Champaign county if not added
	//first find the county id
	chFilter := bson.D{primitive.E{Key: "name", Value: "Champaign"}}
	var champaignCounty *county
	err = counties.FindOne(chFilter, &champaignCounty, nil)
	if err != nil {
		return err
	}
	if champaignCounty == nil {
		return errors.New("there is no a Champaign county")
	}

	//check if added
	filter := bson.D{primitive.E{Key: "app_version", Value: "2.6"}, primitive.E{Key: "county_id", Value: champaignCounty.ID}}
	var items []*model.CRules
	err = cRules.Find(filter, &items, nil)
	if err != nil {
		return err
	}
	if len(items) <= 0 {
		log.Println("there are no symptoms rules for version 2.6 and Champaign county, so we need to add initial data")

		data, err := ioutil.ReadFile("./driven/storage/rules_2.6.json")
		if err != nil {
			return err
		}
		d := model.CRules{AppVersion: "2.6", CountyID: champaignCounty.ID, Data: string(data)}
		_, err = cRules.InsertOne(&d)
		if err != nil {
			return err
		}
	} else {
		log.Println("there are symptoms rules for version 2.6 and Champaign county, so nothing to do")
	}

	log.Println("CRules checks passed")
	return nil
}

func (m *database) applyTraceExposuresChecks(traceExposures *collectionWrapper) error {
	log.Println("apply traceExposures checks.....")

	//add index
	err := traceExposures.AddIndex(bson.D{primitive.E{Key: "date_added", Value: 1}}, false)
	if err != nil {
		return err
	}

	//add index
	err = traceExposures.AddIndex(bson.D{primitive.E{Key: "timestamp", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("traceExposures checks passed")
	return nil
}

func (m *database) applyAccessRulesChecks(accessRules *collectionWrapper) error {
	log.Println("apply accessRules checks.....")

	//add index - unique
	err := accessRules.AddIndex(bson.D{primitive.E{Key: "county_id", Value: 1}}, true)
	if err != nil {
		return err
	}

	log.Println("accessRules checks passed")
	return nil
}

func (m *database) applyUINOverridesChecks(uinoverrides *collectionWrapper) error {
	log.Println("apply uinOverrides checks.....")

	//add index - unique
	err := uinoverrides.AddIndex(bson.D{primitive.E{Key: "uin", Value: 1}}, true)
	if err != nil {
		return err
	}

	//add index
	err = uinoverrides.AddIndex(bson.D{primitive.E{Key: "category", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("uinOverrides checks passed")
	return nil
}

func (m *database) onDataChanged(changeDoc map[string]interface{}) {
	if changeDoc == nil {
		return
	}
	log.Printf("onDataChanged: %+v\n", changeDoc)
	ns := changeDoc["ns"]
	if ns == nil {
		return
	}
	nsMap := ns.(map[string]interface{})
	coll := nsMap["coll"]

	if "configs" == coll {
		log.Println("configs collection changed")

		if m.listener != nil {
			m.listener.OnConfigsChanged()
		}
	} else {
		log.Println("other collection changed")
	}
}
