package xfs

func NewNavigator(fn ...TraverseOptionFn) TraverseNavigator {
	o := composeTraverseOptions(fn...)

	var subject navigatorSubject

	switch o.Subscription {
	case SubscribeAny:
		subject = &universalNavigator{
			navigator: navigator{options: &o},
		}

	case SubscribeFolders:
		subject = &foldersNavigator{
			navigator: navigator{options: &o},
		}

	case SubscribeFiles:
		subject = &filesNavigator{
			navigator: navigator{options: &o},
		}
	}
	nav := navigatorController{
		subject: subject,
	}

	return nav
}
