package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a repreentation of a mongo collection, like an SQL table
type Collection struct {
	col *mongo.Collection
}

// NewCollection creates a new collection
func NewCollection(provideDB DBProviderFunc, cn string) Collection {
	return Collection{
		col: provideDB().Collection(cn),
	}
}

// InsertOneDoc adds a new document to the database
func (c Collection) InsertOneDoc(doc interface{}) (*mongo.InsertOneResult, error) {
	return c.col.InsertOne(context.Background(), doc)
}

// InsertBatch adds a list of documents to the database
func (c *Collection) InsertBatch(doc []interface{}) (*mongo.InsertManyResult, error) {
	return c.col.InsertMany(context.Background(), doc)
}

// FindByID finds by ID, a document in the mongo database
func (c Collection) FindByID(ID string, r interface{}) error {
	filter := bson.M{"id": ID}
	err := c.col.FindOne(context.Background(), filter).Decode(r)
	return err
}

// FindByField finds by a specific field, the FIRST document in the mongo database. See FindLatestByField to get the most recent document
func (c Collection) FindByField(key string, value string, r interface{}) error {
	filter := bson.M{key: value}
	err := c.col.FindOne(context.Background(), filter).Decode(r)
	return err
}

// FindByFilter finds by a passing in a filter
func (c Collection) FindByFilter(filter bson.M, r interface{}) error {
	err := c.col.FindOne(context.Background(), filter).Decode(r)
	return err
}

// FindLatestByField finds by a specific field, the latest document in Mongo
func (c Collection) FindLatestByField(key string, value string, r interface{}) error {
	filter := bson.M{key: value}
	newOpt := options.FindOneOptions{Sort: bson.M{"_id": -1}}
	err := c.col.FindOne(context.Background(), filter, &newOpt).Decode(r)
	return err
}

// FindAll returns all the documents in the collection
// the 'onEach' function is called for each match as the cursor iterates
func (c *Collection) FindAll(onEach func(c *mongo.Cursor) error) error {
	ctx := context.Background()

	filter := bson.M{}
	cur, err := c.col.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		if err := onEach(cur); err != nil {
			return err
		}
	}
	if err := cur.Err(); err != nil {
		return err
	}
	return nil
}

//FindMulti returns multiple documents that match the filter options criteria
// the 'onEach' function is called for each match as the cursor iterates
func (c *Collection) FindMulti(key string, value interface{}, onEach func(c *mongo.Cursor) error) error {
	ctx := context.Background()

	filter := bson.M{key: value}
	cur, err := c.col.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		if err := onEach(cur); err != nil {
			return err
		}
	}
	if err := cur.Err(); err != nil {
		return err
	}
	return nil
}

// FindMultiWithFilter returns multiple documents that match the filter options criteria
// the 'onEach' function is called for each match as the cursor iterates
func (c *Collection) FindMultiWithFilter(filter interface{}, onEach func(c *mongo.Cursor) error) error {
	ctx := context.Background()

	cur, err := c.col.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		if err := onEach(cur); err != nil {
			return err
		}
	}
	if err := cur.Err(); err != nil {
		return err
	}
	return nil
}

// Replace replaces an existing document in the database
func (c Collection) Replace(ID string, replacement interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": ID}
	return c.col.ReplaceOne(context.Background(), filter, replacement)
}

// ReplaceWithFilter replaces an existing document, given a filter, in the database
func (c Collection) ReplaceWithFilter(key, value string, replacement interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{key: value}
	return c.col.ReplaceOne(context.Background(), filter, replacement)
}

// Update updates a specific field in an existing document in the database
func (c Collection) Update(ID string, key string, u interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": ID}
	update := bson.M{"$set": bson.M{key: u}}
	return c.col.UpdateOne(context.Background(), filter, update)
}

// UpdateObject updates existing document in the database
func (c *Collection) UpdateObject(id string, changes interface{}) error {
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", changes},
	}
	_, err := c.col.UpdateOne(context.Background(), filter, update)
	return err
}

// UpdateOneWithFilterOptions updates with filter conditions specified
func (c *Collection) UpdateOneWithFilterOptions(filter interface{}, changes interface{}) error {
	update := bson.D{
		{"$set", changes},
	}
	_, err := c.col.UpdateOne(context.Background(), filter, update)
	return err
}

// Delete deletes a document from a collection in the data
func (c Collection) Delete(ID string) error {
	filter := bson.M{"id": ID}
	_, err := c.col.DeleteOne(context.Background(), filter)
	return err
}

// deleteAll deletes all documents in a collection
func (c *Collection) deleteAll() error {
	_, err := c.col.DeleteMany(context.Background(), bson.D{})
	return err
}

// IsNotFoundError checks if a mongo error is no entity found error
func IsNotFoundError(err error) bool {
	return err == mongo.ErrNoDocuments
}
