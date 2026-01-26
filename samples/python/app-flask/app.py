from flask import Flask, jsonify
import os
import logging

app = Flask(__name__)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

APP_VERSION = os.getenv('APP_VERSION', '1.0.0')
ENVIRONMENT = os.getenv('ENVIRONMENT', 'development')


@app.route('/', methods=['GET'])
def hello():
    return jsonify({
        'message': 'Welcome to Flask App',
        'version': APP_VERSION,
        'environment': ENVIRONMENT,
        'status': 'healthy'
    }), 200


@app.route('/health', methods=['GET'])
def health():
    return jsonify({
        'status': 'healthy',
        'version': APP_VERSION
    }), 200


@app.route('/ready', methods=['GET'])
def readiness():
    return jsonify({
        'ready': True,
        'version': APP_VERSION
    }), 200


@app.route('/api/info', methods=['GET'])
def info():
    return jsonify({
        'app_name': 'Flask-Kubernetes-App',
        'version': APP_VERSION,
        'environment': ENVIRONMENT,
        'description': 'A Flask application deployed on Kubernetes with Helm'
    }), 200


@app.errorhandler(404)
def not_found(error):
    return jsonify({'error': 'Not Found', 'status': 404}), 404


@app.errorhandler(500)
def internal_error(error):
    logger.error(f'Internal server error: {error}')
    return jsonify({'error': 'Internal Server Error', 'status': 500}), 500


if __name__ == '__main__':
    host = os.getenv('FLASK_HOST', '0.0.0.0')
    port = int(os.getenv('FLASK_PORT', 5000))
    debug = os.getenv('FLASK_DEBUG', 'False').lower() == 'true'
    
    logger.info(f'Starting Flask app on {host}:{port}')
    app.run(host=host, port=port, debug=debug)
