import { readFileSync, writeFileSync } from "node:fs";

type SubDivision$1 = {
    id: string;
    name: string;
    state_id: string;
};

type States$1 = {
    [x: string]: SubDivision$1[];
};

type CountryData$1 = {
    name: string;
    native: string;
    phone: string;
    continent: string;
    capital: string;
    currency: string;
    languages: string[];
    emoji: string;
    emojiU: string;
    states: States$1;
};

type CountryList$1 = {
    [x: string]: CountryData$1;
};

export type ConvContinent = {
    id: number;
    name: string;
};

export type ConvCountry = {
    id: number;
    continentId: number;
    name: string;
    code: string;
    phone: string;
    currency: string;
};

export type ConvState = {
    id: number;
    name: string;
    countryId: number;
};

export type ConvCity = {
    id: number;
    name: string;
    stateId: number;
};

export type TransformedData = {
    continents: ConvContinent[];
    countries: ConvCountry[];
    states: ConvState[];
    cities: ConvCity[];
};

const counter = {
    country: 1,
    state: 1,
    cities: 1,
};

const continentMap = {
    Europe: 1,
    "North-America": 2,
    "South-America": 3,
    Africa: 4,
    Asia: 5,
    Oceania: 6,
    Antartica: 7,
};
const countriesArr = [] as ConvCountry[];
const statesArr = [] as ConvState[];
const citiesArr = [] as ConvCity[];

function mapToNameValArr(map: Record<string, unknown>) {
    const out = [];
    for (const k in map) {
        out.push({
            id: map[k],
            name: k,
        });
    }
    return out;
}

export function transformWorldData() {
    const worldStr = readFileSync("./compiled-cities.json", { encoding: "utf-8" });
    const worldData = JSON.parse(worldStr) as CountryList$1;

    for (const key in worldData) {
        const { name, states, continent, currency, phone } = worldData[key];
        const continentId = continentMap[continent as keyof typeof continentMap];
        if (!continentId) {
            console.error("Invalid Continent name", continent);
            throw new Error("invalid continent name!!!");
        }
        const countryId = counter.country++;
        const country = {
            id: countryId,
            code: key,
            continentId,
            phone,
            currency,
            name,
        };
        countriesArr.push(country);

        for (const sKey in states) {
            const subDivision = states[sKey];
            if (!Array.isArray(subDivision)) {
                throw new Error("subDivision for state `" + sKey + "` is not an array!!!");
            }

            const stateId = counter.state++;
            const state = {
                id: stateId,
                countryId,
                name: sKey,
            };
            statesArr.push(state);

            for (const div of subDivision) {
                const { name } = div;
                const cityId = counter.cities++;
                const city = {
                    id: cityId,
                    stateId,
                    name,
                };

                citiesArr.push(city);
            }
        }
    }

    const output = {
        continents: mapToNameValArr(continentMap) as ConvContinent[],
        countries: countriesArr,
        states: statesArr,
        cities: citiesArr,
    };

    writeFileSync("./db-ready-world.json", JSON.stringify(output, null, 2));
    console.log("Done...!");
}

if (import.meta.url === `file://${process.argv[1]}`) {
    transformWorldData();
}
