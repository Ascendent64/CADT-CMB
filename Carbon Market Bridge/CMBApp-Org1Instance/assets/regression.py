import sys
import json
import numpy as np
import statsmodels.api as sm

def perform_regression(tokens, regions):
    try:
        if not isinstance(tokens, list):
            raise ValueError("tokens must be a list")
        if not isinstance(regions, dict):
            raise ValueError("regions must be a dictionary")
        
        X = []
        y = []
        categories = {
            'Registry': {},
            'Location': {},
            'Sector': {}
        }

        for token in tokens:
            try:
                feature_vector = []
                registry = token['unitData']['issuanceData']['projectData']['currentRegistry']
                sector = token['unitData']['issuanceData']['projectData']['sector']
                country = token['unitData']['issuanceData']['projectData']['locations'][0]['country']
                location = None

                for region, countries in regions.items():
                    if country in countries:
                        location = region
                        break

                if registry:
                    if registry not in categories['Registry']:
                        categories['Registry'][registry] = len(categories['Registry'])
                    feature_vector.append(categories['Registry'][registry])

                if location:
                    if location not in categories['Location']:
                        categories['Location'][location] = len(categories['Location'])
                    feature_vector.append(categories['Location'][location])

                if sector:
                    if sector not in categories['Sector']:
                        categories['Sector'][sector] = len(categories['Sector'])
                    feature_vector.append(categories['Sector'][sector])

                if token.get('listedUnitPrice') is not None:
                    X.append(feature_vector)
                    y.append(float(token['listedUnitPrice']))

            except KeyError as ke:
                raise ValueError(f"Missing key in token data: {ke}")
            except Exception as e:
                raise ValueError(f"Error processing token data: {str(e)}")

        if len(X) == 0 or len(y) == 0:
            return {
                'coefficients': {
                    'Registry': {},
                    'Location': {},
                    'Sector': {}
                },
                'baselinePrice': None
            }

        X = np.array(X)
        y = np.array(y)

        X = sm.add_constant(X) 
        model = sm.OLS(y, X)
        results = model.fit()

        coefficients = results.params
        intercept = coefficients[0]

        tokensCoefficientData = {
            'Registry': {},
            'Location': {},
            'Sector': {}
        }

        for category in ['Registry', 'Location', 'Sector']:
            for variant, index in categories[category].items():
                coef = coefficients[index + 1]  
                coef = min(0, coef)  
                tokensCoefficientData[category][variant] = {
                    'coefficient': coef,
                    'total': sum(1 for token in tokens if token.get('listedUnitPrice') is not None and (
                        (category == 'Registry' and token['unitData']['issuanceData']['projectData']['currentRegistry'] == variant) or
                        (category == 'Location' and location == variant) or
                        (category == 'Sector' and token['unitData']['issuanceData']['projectData']['sector'] == variant)
                    ))
                }

        return {
            'coefficients': tokensCoefficientData,
            'baselinePrice': intercept
        }

    except Exception as e:
        raise RuntimeError(f"Error in perform_regression: {str(e)}")

if __name__ == "__main__":
    try:
        input_data = json.loads(sys.stdin.read())
        tokens = input_data['tokens']
        regions = input_data['regions']

        regression_results = perform_regression(tokens, regions)
        print(json.dumps(regression_results))
    except json.JSONDecodeError:
        print("Error: Failed to parse input JSON.", file=sys.stderr)
        sys.exit(1)
    except KeyError as ke:
        print(f"Error: Missing key in input data: {ke}", file=sys.stderr)
        sys.exit(1)
    except ValueError as ve:
        print(f"Error: {str(ve)}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: Unexpected error: {str(e)}", file=sys.stderr)
        sys.exit(1)
