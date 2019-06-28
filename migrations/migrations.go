package migrations

type Migration interface{
	Guard(afero.Fs, string ) (bool, error) //returns true if migration is runnable, otherwise error
	Migrate(afero.Fs, string ) (string, error) //Returns path to config and err
}

func RunMigrations(fs, configFile) err{
	migrations := []migration{}
		for _,migration := range migrations{
			shouldRun, err := migration.Guard(fs,configFile)
			if err != nil{
				return err
			}
			if shouldRun == false{
				continue
			}
			
			configFile, err := migration.Migrate(fs, configFile)
			if err != nil{
				return err
			}
		}
}