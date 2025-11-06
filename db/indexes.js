import { MongoClient } from "mongodb";

async function setupCollectionsAndIndexes(uri) {
  const client = new MongoClient(uri);

  try {
    await client.connect();
    const db = client.db("mileapp_db");

    await db.createCollection("users");
    console.log("created collection: users");

    await db.createCollection("tasks");
    console.log("created collection: tasks");

    const usersCollection = db.collection("users");

    await usersCollection.createIndex({ email: 1 }, { unique: true });
    console.log("created index on users.email (unique)");

    const tasksCollection = db.collection("tasks");

    await tasksCollection.createIndex({ status: 1 });
    console.log("created index on tasks.status");

    await tasksCollection.createIndex({ priority: 1 });
    console.log("created index on tasks.priority");

    await tasksCollection.createIndex({ due_date: 1 });
    console.log("created index on tasks.due_date");

    await tasksCollection.createIndex({ created_at: -1 });
    console.log("created index on tasks.created_at (descending)");

    console.log("collections and indexes setup complete.");
  } catch (err) {
    console.error("failed to setup collections or indexes:", err);
  } finally {
    await client.close();
  }
}

setupCollectionsAndIndexes("mongodb://localhost:27017");
