package resolvers

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
)

func getResolvedPersonFromActor(actor vocab.ActivityStreamsActorProperty) (vocab.ActivityStreamsPerson, error) {
	var err error
	var person vocab.ActivityStreamsPerson

	personCallback := func(c context.Context, p vocab.ActivityStreamsPerson) error {
		person = p
		return nil
	}

	for iter := actor.Begin(); iter != actor.End(); iter = iter.Next() {
		if iter.IsIRI() {
			fmt.Println("Is IRI", iter.GetIRI())
			iri := iter.GetIRI()
			c := context.TODO()
			if e := ResolveIRI(iri.String(), c, personCallback); e != nil {
				err = e
			}
		} else if iter.IsActivityStreamsPerson() {
			p := iter.GetActivityStreamsPerson()
			person = p
		}
	}

	return person, err
}

func getPersonFromFollow(activity vocab.ActivityStreamsFollow, c context.Context) (vocab.ActivityStreamsPerson, error) {
	fmt.Println("getPersonFromFollow...")

	return getResolvedPersonFromActor(activity.GetActivityStreamsActor())
}

func MakeFollowRequest(activity vocab.ActivityStreamsFollow, c context.Context) (*apmodels.ActivityPubActor, error) {
	person, err := getPersonFromFollow(activity, c)
	if err != nil {
		return nil, err
	}

	fmt.Println(activity.GetJSONLDId().Get().String())

	followRequest := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
	}

	return &followRequest, nil
}

func MakeUnFollowRequest(activity vocab.ActivityStreamsUndo, c context.Context) *apmodels.ActivityPubActor {
	person, err := getResolvedPersonFromActor(activity.GetActivityStreamsActor())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if person == nil {
		return nil
	}

	fmt.Println(activity.GetJSONLDId().Get().String())

	request := apmodels.ActivityPubActor{
		ActorIri:  person.GetJSONLDId().Get(),
		FollowIri: activity.GetJSONLDId().Get(),
		Inbox:     person.GetActivityStreamsInbox().GetIRI(),
	}

	return &request
}