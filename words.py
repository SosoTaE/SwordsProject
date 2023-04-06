import pymongo

_str = "mongodb+srv://sosotae:SosoTaENaoko;@sword.yah6gsy.mongodb.net/?retryWrites=true&w=majority"
__str = "mongodb://localhost:27017"

client = pymongo.MongoClient(__str)

data = list(client.texts.words.find())

client = pymongo.MongoClient(_str)

client.texts.words.insert_many([{"word":doc["word"]} for doc in data])

