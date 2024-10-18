import { PathConstants } from "./pathConstants";
import WeatherForm from "../components/WeatherForm/WeatherForm";
import PageNotFound from "../components/PageNotFound/PageNotFound";

export const routes = [
  {
    path: PathConstants.HOME,
    element: <WeatherForm />,
    errorElement: <PageNotFound />,
  },
];
