import { PathConstants } from "@routes/pathConstants"
import WeatherForm from "@components/WeatherForm/WeatherForm"
import PageNotFound from "@components/PageNotFound/PageNotFound"

export const routes = [
  {
    path: PathConstants.HOME,
    element: <WeatherForm />,
    errorElement: <PageNotFound />,
  },
]
