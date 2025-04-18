import PageNotFound from "@components/PageNotFound/PageNotFound"
import WeatherForm from "@components/WeatherForm/WeatherForm"
import { PathConstants } from "@routes/pathConstants"

export const routes = [
  {
    path: PathConstants.HOME,
    element: <WeatherForm />,
    errorElement: <PageNotFound />,
  },
]
