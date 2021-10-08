import requests

class Test:
	@staticmethod
	def upload(file: str):
		url = "http://127.0.0.1:9779/api/upload"
		files = {'file': open(file, 'rb')}
		res = requests.post(url, files = files)
		print(res.json())


if __name__ == '__main__':
	Test.upload('H:\\code\\images\\375.jpg')
