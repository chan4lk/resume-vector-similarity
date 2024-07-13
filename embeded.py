from flask import Flask, request, jsonify
from waitress import serve
from sentence_transformers import SentenceTransformer
model = SentenceTransformer('/Users/chandima/repos/all-mpnet-base-v2')

mbed = Flask(__name__)
@mbed.route('/api/embeddings', methods=['POST'])
def run():
    data = request.json
    embedding = model.encode(data['prompt'])
    return jsonify({'embedding': embedding.tolist()}), 200, {'Content-Type': 'application/json; charset=utf-8'}

if __name__ == '__main__':
    serve(mbed, port=11333, threads=16)