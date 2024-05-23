from flask import Flask, render_template, request
import requests
from bs4 import BeautifulSoup
import pandas as pd

app = Flask(__name__, static_url_path='/static')

@app.route("/")
def get_data():
    return render_template('Dashboard.html') 

def company():
    compet_url = 'https://en.wikipedia.org/wiki/List_of_unicorn_startup_companies#cite_note-16'
    response = requests.get(compet_url)
    soup = BeautifulSoup(response.text, 'html.parser')
    tables = soup.findAll("table", { "class" : "wikitable" })
    table1 = tables[2]
    df = pd.read_html(str(table1))[0]

    # Rename columns as needed
    df.columns = ['Company', 'Valuation(US$ billions)', 'Valuation date', 'Industry', 'Country', 'Founders']
    df.to_csv('output.csv', index=False)
    return df

@app.route("/Companies.html")
def show_companies():
    df = company()
    return render_template("Companies.html", companies=df.to_dict(orient='records'))

def get_filtered_data(df, industry=None, country=None):
    # Fill NaN values with an empty string
    df = df.fillna('')

    # Print column names for debugging
    print("Column Names:", df.columns)

    # Filtering data based on industry and country
    if industry:
        df = df[df['Industry'].str.contains(industry, case=False)]
    if country:
        df = df[df['Country'].str.contains(country, case=False)]

    return df.to_dict(orient='records')

@app.route("/Compete.html", methods=['GET', 'POST'])
def compete():
    if request.method == 'POST':
        industry = request.form['industry']
        country = request.form['country']
        df = company()
        filtered_data = get_filtered_data(df, industry=industry, country=country)
        return render_template('Compete.html', filtered_data=filtered_data)
    else:
        return render_template('Compete.html', filtered_data=None)

if __name__ == '__main__':
    app.run(debug=True)
