Server Side Api (ver0.1 Draft)
=============================
## 共有Error
	error == "": 成功.
	error == "err_param": 参数错误.
	error == "err_auth": 身份验证错误，没登陆或者token超时.
	error == "err_internal": 服务器内部错误.
	
## 所有的whapi下的接口都需要带上token参数，比如
- 		{ServerURL}/whapi/player/getinfo?token=16995a9581c74b18ad1584ad9c68d245
- 		{ServerURL}/whapi/player/create?warlord=3&token=16995a9581c74b18ad1584ad9c68d245

## 测试server: 42.121.107.155

--------------------------------------------
## authapi/register
- Method: GET
- Desc: 注册帐号
- Param:
-		username = {STRING(40)}
		password = {STRING(40)}
- Example: 
- 		{ServerURL}/authapi/register?username=foo&password=bar123
- Return:
-		{
			error: {STRING}	/* 错误 */
			token: {STRING}	/* 用户token，用于后续会话的身份验证。（同名cookie里有同样的信息） */
		}
		error == "err_account_exist": 用户名已存在.

## authapi/login
- Method: GET
- Desc: 用户登入
- Param:
-		username = {STRING(40)}
		password = {STRING(40)}
- Example: 
- 		{ServerURL}/authapi/login?username=foo&password=bar123
- Return:
-		{
			error: {STRING},	/* 错误 */
			token: {STRING}	/* 用户token，用于后续会话的身份验证。（同名cookie里有同样的信息） */
			playerExist: {BOOL} /* 用户的角色是否已经创建 */
		}
		error == "err_not_match": 用户名与密码不匹配，即用户名未注册或者密码错误.

## whapi/version
- Method: GET
- Desc: 获得当前版本信息
- Param:
-		None
- Example: 
- 		{ServerURL}/whapi/version
- Return:
-		{
			version: [{INT}, {INT}, {INT}],		/* 当前服务器版本号，相当于x.x.x */
			minVersion: [{INT}, {INT}, {INT}]	/* 最小兼容版本号, 客户端必须大于此版本号才可运行 */
		}
		
## whapi/player/create
- Method: GET
- Desc: 选主角，初始化玩家信息
- Param: 
- 		warlord = {INT}		/* ProtoID */ 
- Example:
- 		{ServerURL}/whapi/player/create?warlord=3&token=16995a9581c74b18ad1584ad9c68d245
- Return: 
- 		{
			error: {STRING},    	/* 错误 */
			warlord: {
				id: {INT},			/* 实体id（流水号） */
				protoId: {INT},		/* 原型id（卡片类型） */
				level: {INT},		/* 等级 */
				exp: {INT},			/* 经验值 */
				hp: {INT},			/* 基础hp */
				atk: {INT},			/* 基础atk */
				def: {INT},			/* 基础def */
				wis: {INT},			/* 基础wis */
				agi: {INT},			/* 基础agi */
				hpCrystal: {INT},	/* 水晶增加的hp */
				atkCrystal: {INT},	/* 水晶增加的atk */
				defCrystal: {INT},	/* 水晶增加的def */
				wisCrystal: {INT},	/* 水晶增加的wis */
				agiCrystal: {INT},	/* 水晶增加的agi */
				hpExtra: {INT},		/* 进化增加的hp */
				atkExtra: {INT},	/* 进化增加的atk */
				defExtra: {INT},	/* 进化增加的def */
				wisExtra: {INT},	/* 进化增加的wis */
				agiExtra: {INT},	/* 进化增加的agi */
				skill1Id: {INT},	/* 技能1 id */
				skill1Level: {INT},	/* 技能1 level */
				skill1Exp: {INT},	/* 技能1 exp */
				skill2Id: {INT},	/* 技能2 id */
				skill2Level: {INT},	/* 技能2 level */
				skill2Exp: {INT},	/* 技能2 exp */
				skill3Id: {INT},	/* 技能3 id */
				skill3Level: {INT},	/* 技能3 level */
				skill3Exp: {INT},	/* 技能3 exp */
			}
		}
		error == "err_player_exist": player已存在
		error == "err_create_warlord": warlord创建失败

## whapi/player/getinfo
- Method: GET
- Desc: 获取当前登入玩家的信息
- Param: 无
- Example: 
- 		{ServerURL}/whapi/player/getinfo?token=16995a9581c74b18ad1584ad9c68d245
- Return: 
-		{
			error: {STRING},    	/* 错误 */
			userId: {INT},			/* 用户ID */
			name: {STRING},			/* 用户名称 */
			warload: {INT},			/* 主角卡ID */
			money: {INT},			/* 钱 */
			inZoneId: {INT},    	/* 如果已进入zone返回zoneId，否则返回0 */
			lastZoneId: {INT}, 		/* 已解锁的地图的最后ID */
			ap: {INT},				/* 当前行动值，用于推图 */
			maxAp: {INT},			/* 最大行动值 */
			apAddRemain: {INT},		/* 距下次增加ap时间（秒） */
			xp: {INT},				/* 当前活动行动值，用于pvp，boss战等 */
			maxXp: {INT},			/* 最大活动行动值 */
			xpAddRemain: {INT},		/* 距下次增加xp时间（秒） */
			lastFormation: {INT},	/* 最新得到的阵形 */
			currentBand: {INT},		/* 当前选择band index, 从0开始*/
			maxCardNum: {INT},		/* 最大卡片数量 */
			wagonGeneral: {INT},	/* general仓库中的货物总数 */
 			wagonTemp: {INT},		/* temp仓库中的货物总数 */
 			wagonSocial: {INT},		/* social仓库中的货物总数 */
			pvpWinStreak: {INT}		/* 连胜数*/
			cards: 	[				/* 当前用户拥有的卡片信息，为cardEntity对象的数组*/
						{
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
						},
						... 
			], 
			bands: 	[							/* band信息 *
						{
							index: {INT},		/* band index,区间为[0,2] */
							formation: {INT}, 	/* 阵形id */
							members: [{INT or null}, ...]	/* band 成员，使用cardEntityId.空位用null.数量必须与阵形匹配 */
						},
						...
			],		
			items:	[							/* item信息 */
						{ 
							id: {INT},		/* item id */
							num: {INT} 		/* item 数量 */
						},		
						...
			]
		}
		error == "err_player_not_exist": player未创建
		
## whapi/player/setband
- Method: POST
- Desc: 设置band
- Input:
- 		[
			{
				index: {INT},		/* band index,区间为[0,2] */
				formation: {INT}, 	/* 阵形id */
				members: [{INT or null}, ...]	/* 成员，使用cardEntityId.空位用null.数量必须与阵形匹配 */
			}
			...
		]
- Example: 
-		[
			{"index":0, "formation":23, "members":[34, 643, null, 456, null, 54]},
			{"index":2, "formation":25, "members":[454, 543, 43, 37, 77, 54]}
		]
- Return: 
- 		{
 			error: {STRING},   	/* 错误 */
		}
		

## whapi/player/useitem
- Method: POST
- Desc: 使用物品。注意：如果物品的使用对象是全体，请勿指定targets（留空）。
- Input:
-		{
			itemid: {INT}			/* 使用物品id */
			targets: [{INT}, ...] 	/* 物品使用对象cardEntityId, 如果没有指定对象可以留空 */
			allout: {BOOL}			/* 只在使用小号角时有用。指定是否用小号角加满xp */
		}
- Example: 
-		{ServerURL}/whapi/player/useitem?token=16995a9581c74b18ad1584ad9c68d245
		{
			"itemid": 2,
			"targets": [233, 56334, 5432]
		}
- Return: 
- 		{
 			error: {STRING},   	/* 错误 */
			itemId: {INT},		/* 被使用物品 */
			itemNum: {INT},		/* 物品剩余数量 */
		}

## whapi/player/time
- Method: Get
- Desc: 获取Player相关的时间信息
- Param: -
- Example: {ServerURL}/whapi/player/gettime?token=16995a9581c74b18ad1584ad9c68d245
- Return: 
	{

		ap: {INT},				/* 当前行动值，用于推图 */
		apAddRemain: {INT},		/* 距下次增加ap时间（秒） */
		xp: {INT},				/* 当前活动行动值，用于pvp，boss战等 */
		xpAddRemain: {INT},		/* 距下次增加xp时间（秒） */

		/*TODO: PVP Time */
	}

## whapi/zone/enter
- Method: GET
- Desc: 获取当前地图信息
- Param: 
- 		zoneid = {INT}		/* 地图id *
		bandidx = {INT}		/* 指定使用第几个band， 从0开始 */
- Example: {ServerURL}/whapi/zone/enter?zoneid=10&bandidx=0&token=16995a9581c74b18ad1584ad9c68d245
- Return: 
- 		{
 			error: {STRING},   				/* 错误 */
			zoneId: {INT}, 					/* 当前的Zone ID */
			startPos: {x:{INT}, y:{INT}},	/* 起点位置 */
			goalPos: {x:{INT}, y:{INT}},	/* 终点位置 */
			currPos: {x:{INT}, y:{INT}},	/* 当前位置 */
			redCase: {INT}					/* 红宝箱获得数量 */
			goldCase: {INT}					/* 金宝箱获得数量 */
			objs: [
				[
				
					{INT}, 		/* x坐标 */
				 	{INT}, 		/* y坐标 */
					{INT}		/* 类型 */
									/* 1:木箱 */
									/* 2:红宝箱 */
									/* 3:金宝箱 */
									/* 4:小钱袋 */
									/* 5:大钱袋 */
									/* 6:pvp */
									/* <0:绝对值为Monster Group的id */
				], 
				...
			],
			events:[
				{x: {INT}, 		/* x坐标 */
				 	y: {INT}, 		/* y坐标 */
				 	startDialog: {INT},
				 	endDialog: {INT},
				 	monsterId: {INT}
				}, 
				...
			],
			band: {
				formation: {INT}	/* 阵形id */
				members: [
					{
						id: {INT}, 
						hp: {INT}
					}, /* id为CardEntityId, hp为该Card的当前血量*/ 
					...
				]
			},
			enterDialogue: {INT},
			completeDialogue: {INT}
		}


## whapi/zone/get
- Method: GET
- Desc: 获取当前进入的地图的游戏数据
- Input: 无
- Example: {ServerURL}/whapi/zone/get?token=16995a9581c74b18ad1584ad9c68d245
- Return: 
- 		{
 			error: {STRING},   	/* 错误 */
			zoneId: {INT}, 		/* 当前的Zone ID */
			startPos: {x:{INT}, y:{INT}},	/* 起点位置 */
			goalPos: {x:{INT}, y:{INT}},	/* 终点位置 */
			currPos: {x:{INT}, y:{INT}},	/* 当前位置 */
			redCase: {int}					/* 红宝箱获得数量 */
			goldCase: {int}					/* 金宝箱获得数量 */
			objs: [
				[
				
					{INT}, 		/* x坐标 */
				 	{INT}, 		/* y坐标 */
					{INT}		/* 类型 */
									/* 1:木箱 */
									/* 2:红宝箱 */
									/* 3:金宝箱 */
									/* 4:小钱袋 */
									/* 5:大钱袋 */
									/* 6:pvp */
									/* <0:绝对值为Monster Group的id */
				], 
				...
			],
			events:[{x: {INT}, 		/* x坐标 */
				 	y: {INT}, 		/* y坐标 */
				 	startDialog: {INT},
				 	endDialog: {INT},
				 	monsterId: {INT}}, ...
					],
			band: {
				formation: {INT}	/* 阵形id */
				members: [
					{
						id: {INT}, 
						hp: {INT}
					}, /* id为CardEntityId, hp为该Card的当前血量*/ 
					...
				]
			},
			enterDialogue: {INT},
			completeDialogue: {INT}
		}

## whapi/zone/withdraw
- Method: GET
- Desc: 退出zone
- Param: 无
- Example: {ServerURL}/whapi/zone/withdraw?token=16995a9581c74b18ad1584ad9c68d245
- Return: 
- 		{
 			error: {STRING},   	/* 错误 */
		}

## whapi/zone/move
- Method: POST
- Desc: 在地图中移动
- Param:
-		[					/* 顺序走过的所有坐标 */
			[INT, INT]，		/* xy坐标 */
			[INT, INT]			
		]
- Example: {ServerURL}/whapi/zone/move?x=10&y=17
- Return: 
- 		{
			error: {STRING},		/* 错误 */
			currPos: {x:{INT}, y:{INT}},/* 地图中的当前位置 */
			ap: {INT},				/* 当前行动值 */
			nextAddApTime: {INT},	/* 距下次加ap还有几秒 */
			redCaseAdd: {INT},		/* 获得的红宝箱数量（0或1） */
			goldCaseAdd: {INT},		/* 获得的金宝箱数量（0或1） */
			items: [				/* 获得item,没有的话为空数组[] */
				{
					id: {INT},		/* 获得item类型 */
					num: {INT}		/* 获得item数量 */
				},
				...
			],
			cards: [					/* 奖励的卡片，没有的话，数组长度为0 */
				{
					id: {INT},			/* 实体id（流水号） */
					protoId: {INT},		/* 原型id（卡片类型） */
					level: {INT},		/* 等级 */
					exp: {INT},			/* 经验值 */
					hp: {INT},			/* 基础hp */
					atk: {INT},			/* 基础atk */
					def: {INT},			/* 基础def */
					wis: {INT},			/* 基础wis */
					agi: {INT},			/* 基础agi */
					hpCrystal: {INT},	/* 水晶增加的hp */
					atkCrystal: {INT},	/* 水晶增加的atk */
					defCrystal: {INT},	/* 水晶增加的def */
					wisCrystal: {INT},	/* 水晶增加的wis */
					agiCrystal: {INT},	/* 水晶增加的agi */
					hpExtra: {INT},		/* 进化增加的hp */
					atkExtra: {INT},	/* 进化增加的atk */
					defExtra: {INT},	/* 进化增加的def */
					wisExtra: {INT},	/* 进化增加的wis */
					agiExtra: {INT},	/* 进化增加的agi */
					skill1Id: {INT},	/* 技能1 id */
					skill1Level: {INT},	/* 技能1 level */
					skill1Exp: {INT},	/* 技能1 exp */
					skill2Id: {INT},	/* 技能2 id */
					skill2Level: {INT},	/* 技能2 level */
					skill2Exp: {INT},	/* 技能2 exp */
					skill3Id: {INT},	/* 技能3 id */
					skill3Level: {INT},	/* 技能3 level */
					skill3Exp: {INT},	/* 技能3 exp */
				},
				... 
					],
			
			eventId: {INT}, 				/* 当前格子的event id */
			catchMons: [{INT}, ...]			/* 可捕捉怪物prototypeId */
			pvpBands: [						/* pvp匹配到的3个队伍，如果没有为[] */
				{							/* pvp队伍信息 */
					userId:{INT}			/* 对手用户id */
					userName:{STRING}		/* 对手用户名 */
					userLevel: {INT}        /* 对手的等级 */
					formation:{INT}			/* 阵形 */
					cards:[					/* 对手band */
						{					/* band中的成员，如果位置为空，则为null */
						    protoId:{INT},			/* 原型id（卡片类型） */
						    level:{INT},			/* 等级 */
						    hp:{INT},				/* hp总和 */
						    atk:{INT},				/* atk总和 */
						    def:{INT},				/* def总和 */
						    wis:{INT},				/* wis总和 */
						    agi:{INT},				/* agi总和 */
						    skill1Id:{INT},			/* 技能1 id */
						    skill1Level:{INT},		/* 技能1 等级 */
						    skill2Id:{INT},			/* 技能2 id */
						    skill2Level:{INT},		/* 技能2 等级 */
						    skill3Id:{INT},			/* 技能3 id */
						    skill3Level:{INT}		/* 技能3 等级 */
						},
						...
					]
				},
				...
			]
		}

## whapi/zone/battleresult
- Method: POST
- Desc: 获得Battle 的结果
- Input: 
- 		{
			isWin: {BOOL}, 		/* 战斗是否胜利 */
			members: [
				{
					id: {INT},   /* CardEnity Id */
					hp: {INT}
				}, /* id为CardEntityId, hp为该Card的当前血量*/ 
				...
			],
		}
- Example: 
		{ServerURL}/whapi/zone/battleresult?token=16995a9581c74b18ad1584ad9c68d245
		
		{"isWin":true,"members":[null,{"hp":1743,"id":119},null]}
		
- Return: 
- 		{
			error : {STRING}  		/* 通用格式，如果为空字符串，则为正确 */
			members: [				/* 成员更新经验值 */
				{
					id: {INT},		/* card entity id */
					exp: {INT}，		/* 更新后的经验值 */
				},
				...
			],
			levelups: [				/* 升级了的成员 */
				{
					id: {INT}, 		/* card entity id */
					level: {INT}，	/* 更新后的等级 */
					hp: {INT},		/* 更新后基础hp */
					atk: {INT},		/* 更新后基础atk */
					def: {INT},		/* 更新后基础def */
					wis: {INT},		/* 更新后基础wis */
					agi: {INT}		/* 更新后基础agi */
				}
			],
			cards: [					/* 奖励的卡片，没有的话，数组长度为0 */
						{
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
						},
						... 
					],
			items: [				/* 过关奖励的道具，没有的话为空数组[] */
				{					
					id: {INT},		/* 获得item类型 */
					num: {INT}		/* 获得item数量 */
				},
				...
			]
		}

## whapi/zone/catchmonster
- Method: POST
- Desc: 抓怪结果
- Input: 
- 		{
			catchItem: {INT}	/* 5：大抓怪药，6：小抓怪药 */
		}
- Example: 
		{ServerURL}/whapi/zone/catchmonster?token=16995a9581c74b18ad1584ad9c68d245
		
		{"catchItem":6}
		
- Return: 
- 		{
			error : {STRING}  		/* 通用格式，如果为空字符串，则为正确 */
			catchedMons: [				/* 当前用户拥有的卡片信息，为cardEntity对象的数组*/
						{
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
						},
						... 
			]，
			powder: {INT},		/*剩余小抓怪药*/
			advPowder: {INT}	/*剩余大抓怪药*/
		}

## whapi/zone/complete
- Method: GET
- Desc: 完成当前地图
- Param: 无
- Example: {ServerURL}/whapi/zone/complete
- Return: 
- 		{
			error: {STRING}  		/* 通用格式，如果为空字符串，则为正确 */
			redCase: [				/* 开红宝箱获得item,没有的话为空的object（{}） */
				{
					id: {INT},		/* 获得item类型 */
					num: {INT}		/* 获得item数量 */
				},
				...
			],
			goldCase: [				/* 开金宝箱获得item,没有的话为空的object（{}） */
				{
					id: {INT},		/* 获得item类型 */
					num: {INT}		/* 获得item数量 */
				},
				...
			],
			lastZoneId: {INT}, 			/* 最新解锁的地图ID, 没有为null */
			cards: [					/* 过关奖励的卡片，没有的话，数组长度为0 */
						{
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
						},
						... 
					],
			formation: {INT},		/* 过关奖励的阵型，没有为null */
			maxCardNum: {INT},		/* 过关奖励的玩家卡片数量上限，没有为null */
			maxTradeNum: {INT},		/* 过关奖励的玩家每日交易，没有为null */
			items: [				/* 过关奖励的道具，没有的话为空数组[] */
				{					
					id: {INT},		/* 获得item类型 */
					num: {INT}		/* 获得item数量 */
				},
				...
			]
		}

## whapi/card/getpact
- Method: GET
- Desc: 抽卡
- Param:
		packid = {INT} 		/* 卡包的id */
		num = {INT}			/* 抽卡次数，只有当消耗物品不为whCoin才起作用, 目前上限暂定为10 */
- Example: {ServerURL}/whapi/card/getpact?packid=7&num=5&token=16995a9581c74b18ad1584ad9c68d245
- Return: 
- 		{
 			error: {STRING},   	/* 错误 */
			cards: 	[				/* 当前用户拥有的卡片信息，为cardEntity对象的数组*/
						{
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
						},
						... 
			]
		}

## whapi/card/sell
- Method: POST
- Desc: 卖卡给系统，得到money
- Input:
		[
			{INT},	/* 卡的entity id */
			...
		]

- Example: {ServerURL}/whapi/card/sell?token=16995a9581c74b18ad1584ad9c68d245
		[345, 66, 7456, 74567, 45235]
- Return: 
- 		{
 			error: {STRING},   		/* 错误 */
			cardIds: [{INT}, ...],	/* 卡片id数组，同post数据 */
			money:{INT},			/* 最新money数量 */
			moneyAdd: {INT},		/* 增加的money数量 */
		}

## whapi/card/evolution
- Method: GET
- Desc: 进化
- Param:
		cardid1 = {INT}		/* 主卡id */
		cardid2 = {INT}		/* 副卡id */

- Example: {ServerURL}/whapi/card/evolution?cardid1=3455&cardid2=45678&token=16995a9581c74b18ad1584ad9c68d245
		
- Return: 
- 		{
 			error: {STRING},   		/* 错误 */
			money:{INT},			/* 最新money数量 */
			delCardId: {INT},		/* 副卡ID（被吃的卡） */
			evoCard: {				/* 进化后的主卡 */
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
			}
		}

## whapi/card/sacrifice
- Method: POST
- Desc: 献祭，加技能属性
- Input:
		{
			master: {INT}, /* 主卡 */
			sacrificers: [{INT}, ...] /* 用来献祭的卡 */
		}	

- Example: {ServerURL}/whapi/card/sacrifice?token=16995a9581c74b18ad1584ad9c68d245
		{
			master: 12345,
			sacrificers: [34444, 5456, 74567, 8768, 32534]
		}	

- Return: 
- 		{
 			error: {STRING},   			/* 错误 */
			money:{INT},				/* 最新money数量 */
			sacrificers: [{INT}, ...]	/* 献祭掉的卡的id数组 */
			master: {					/* 加了技能点后的主卡 */ fixme；不需要这么多信息
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
			}
		}

## whapi/card/addcrystal
- Method: POST
- Desc: 卡牌加水晶
- Input:
		{
			card: {INT},	/* 卡牌的CardEntity ID */
			crystal: { 		/* 吃水晶 */
				HP: {INT}, 	/* 要吃的HP水晶的数量 */
				ATK: {INT},	/* 要吃的ATK水晶的数量 */
				DEF: {INT},	/* 要吃的DEF水晶的数量 */
				WIS: {INT},	/* 要吃的WIS水晶的数量 */
				AGI: {INT},	/* 要吃的AGI水晶的数量 */
				GOD: {INT}	/* 暂时不做 */
			} 
		}	

- Example: {ServerURL}/whapi/card/addcrystal?token=16995a9581c74b18ad1584ad9c68d245
		{
			card: 12345,
			crystal:{ 		
				HP: 1,
				ATK: 5,
				DEF: 0,
				WIS: 2,
				AGI: 0,
				GOD: 0
			} 
		}	

- Return: 
- 		{
 			error: {STRING},   		/* 错误 */
 			itemCrystal:{			/* 当前剩余的水晶数量 */
				HP: {INT},			/* 当前剩余的HP水晶数量 */
				ATK: {INT},			/* 当前剩余的ATK水晶数量 */
				DEF: {INT},			/* 当前剩余的DEF水晶数量 */
				WIS: {INT},			/* 当前剩余的WIS水晶数量 */
				AGI: {INT},			/* 当前剩余的AGI水晶数量 */
				GOD: {INT}
 			}
			card: {					/* 吃了水晶后的卡牌信息 */
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
			}
		}		

## whapi/wagon/getcount
- Method: Get
- Desc: 列出wagon中的物品数量
- Param:
- Example: 
-		{ServerURL}/whapi/wagon/getcount?token=16995a9581c74b18ad1584ad9c68d245
- Return: 
- 		{
 			error: {STRING},   			/* 错误 */
 			general: {INT},				/* general仓库中的货物总数 */
 			temp: {INT},				/* temp仓库中的货物总数 */
 			social: {INT}				/* social仓库中的货物总数 */
 		}

## whapi/wagon/list
- Method: POST
- Desc: 列出wagon中的物品
- Param:{
-		wagonIdx: {INT},					/* 0:general 1:temp 2:social*/
		startIdx: {INT},
		count: {INT}  						/* default = 10 */
		}
- Example: 
-		{ServerURL}/whapi/wagon/list?wagonidx=1&token=16995a9581c74b18ad1584ad9c68d245
-		{"wagonIdx":1, "startIdx":10, "count":10}

- Return: 
- 		{
 			error: {STRING},   			/* 错误 */
			wagonIdx: {INT},			/* 0:general 1:temp 2:social */
			totalCount: {INT},			/* 当前wagon idx 下的所有货物总数 */
			items:[						/* 这是item列表 */
				{						
					key: {INT},			/* key */
					itemId: {INT},		/* 物品id */
					itemNum: {INT},		/* 物品数量 */
					desc: {STRING},		/* 项目描述 */
					time: {STRING}		/* 放入时间 */
				}
				...
			],
			cards:[
				{						/* 这是card列表 */
					key: {INT},			/* key */
					cardEntity: {INT},	/* 卡片实例id */
					cardProto: {INT},	/* 卡片原型id */
					desc: {STRING},		/* 项目描述 */
					time: {STRING}		/* 放入时间 */
				},
				...
			]
		}

## whapi/wagon/accept
- Method: POST
- Desc: 取出物品
- Param:
-		{
			keys: [{INT}, {INT}...]	/* 指定要取出的物品的Key */
		}
- Example: 
-		{ServerURL}/whapi/wagon/accept?token=16995a9581c74b18ad1584ad9c68d245
_		{"keys":[234, 545, 2456]}
- Return:
- 		{
 			error: {STRING},   			/* 错误 */
			info: {STRING},				/* "card_full": 卡片已满，只接受了部分 */
										/* "expired": 有物品已过期 */
			generalCount: {INT},		/* general仓库中的货物总数 */
 			tempCount: {INT},			/* temp仓库中的货物总数 */
 			socialCount: {INT}			/* social仓库中的货物总数 */
			acceptedKeys: [{INT}, {INT}, ...] /* 已接收item或card的数组 */
			items: [
				{
					id: {INT},	 		/* item类型 */
					num: {INT}			/* item数量 */
				},
				...
			]
			cards: [
				{					
					id: {INT},			/* 实体id（流水号） */
					protoId: {INT},		/* 原型id（卡片类型） */
					level: {INT},		/* 等级 */
					exp: {INT},			/* 经验值 */
					hp: {INT},			/* 基础hp */
					atk: {INT},			/* 基础atk */
					def: {INT},			/* 基础def */
					wis: {INT},			/* 基础wis */
					agi: {INT},			/* 基础agi */
					hpCrystal: {INT},	/* 水晶增加的hp */
					atkCrystal: {INT},	/* 水晶增加的atk */
					defCrystal: {INT},	/* 水晶增加的def */
					wisCrystal: {INT},	/* 水晶增加的wis */
					agiCrystal: {INT},	/* 水晶增加的agi */
					hpExtra: {INT},		/* 进化增加的hp */
					atkExtra: {INT},	/* 进化增加的atk */
					defExtra: {INT},	/* 进化增加的def */
					wisExtra: {INT},	/* 进化增加的wis */
					agiExtra: {INT},	/* 进化增加的agi */
					skill1Id: {INT},	/* 技能1 id */
					skill1Level: {INT},	/* 技能1 level */
					skill1Exp: {INT},	/* 技能1 exp */
					skill2Id: {INT},	/* 技能2 id */
					skill2Level: {INT},	/* 技能2 level */
					skill2Exp: {INT},	/* 技能2 exp */
					skill3Id: {INT},	/* 技能3 id */
					skill3Level: {INT},	/* 技能3 level */
					skill3Exp: {INT},	/* 技能3 exp */
				},
				...
			]
		}

## whapi/wagon/sellall
- Method: GET
- Desc: 取出物品
- Param:
-		wagonidx = {INT}			/* 0:general 1:temp 2:social */
- Example: 
-		{ServerURL}/whapi/wagon/sellall?wagonidx=1&token=16995a9581c74b18ad1584ad9c68d245
- Return:
- 		{
 			error: {STRING},   			/* 错误 */
			wagonIdx: {INT},			/* 0:general 1:temp 2:social */
			delKeys: [INT, ...],		/* 物品的key，用来找到并删除*/
		}

## whapi/pvp/battleresult
- Method: POST
- Desc: 提交pvp结果
- Input: 
		{
			isWin: {BOOL}, 						/* 战斗是否胜利 */
			foeUserId: {INT},					/* 对手的userId */
			bandIndex: {INT},					/* 玩家出战的band索引(0,1,2) */
			members: [{INT}, {INT}, null,...],	/* 玩家的卡牌成员数组。使用cardEntityId，空位用null。注意只要求提供前排 */
			allout: {BOOL}						/* 是否全力出击 */
			useItem: {INT}						/* 0：不使用 1：小号角 2：大号角 */
		}

- Example: 
-		{ServerURL}/whapi/pvp/battleresult?win=true&token=16995a9581c74b18ad1584ad9c68d245

- Return: 
- 		{
			error : {STRING}  		/* err_timeout: 超时，pvp匹配请求之后太久没战斗 */
			winStreak: {INT},		/* 连胜数 */
			xp: {INT},				/* 当前xp */
			nextAddXpTime: {INT},	/* 距下次加xp还有几秒 */
			smallXpItemNum: {INT},	/* 小号角数量 */
			bigXpItemNum: {INT},	/* 大号角数量 */
			members: [				/* 成员更新经验值 */
				{
					id: {INT},		/* card entity id */
					exp: {INT}，		/* 更新后的经验值 */
				},
				...
			],
			levelups: [				/* 升级了的成员 */
				{
					id: {INT}, 		/* card entity id */
					level: {INT}，	/* 更新后的等级 */
					hp: {INT},		/* 更新后基础hp */
					atk: {INT},		/* 更新后基础atk */
					def: {INT},		/* 更新后基础def */
					wis: {INT},		/* 更新后基础wis */
					agi: {INT}		/* 更新后基础agi */
				}
			],
			cards: [							/* 奖励的卡片，没有的话为空数组 */
						{
							id: {INT},			/* 实体id（流水号） */
							protoId: {INT},		/* 原型id（卡片类型） */
							level: {INT},		/* 等级 */
							exp: {INT},			/* 经验值 */
							hp: {INT},			/* 基础hp */
							atk: {INT},			/* 基础atk */
							def: {INT},			/* 基础def */
							wis: {INT},			/* 基础wis */
							agi: {INT},			/* 基础agi */
							hpCrystal: {INT},	/* 水晶增加的hp */
							atkCrystal: {INT},	/* 水晶增加的atk */
							defCrystal: {INT},	/* 水晶增加的def */
							wisCrystal: {INT},	/* 水晶增加的wis */
							agiCrystal: {INT},	/* 水晶增加的agi */
							hpExtra: {INT},		/* 进化增加的hp */
							atkExtra: {INT},	/* 进化增加的atk */
							defExtra: {INT},	/* 进化增加的def */
							wisExtra: {INT},	/* 进化增加的wis */
							agiExtra: {INT},	/* 进化增加的agi */
							skill1Id: {INT},	/* 技能1 id */
							skill1Level: {INT},	/* 技能1 level */
							skill1Exp: {INT},	/* 技能1 exp */
							skill2Id: {INT},	/* 技能2 id */
							skill2Level: {INT},	/* 技能2 level */
							skill2Exp: {INT},	/* 技能2 exp */
							skill3Id: {INT},	/* 技能3 id */
							skill3Level: {INT},	/* 技能3 level */
							skill3Exp: {INT},	/* 技能3 exp */
						},
						... 
					],
			items: [				/* 奖励的道具，没有的话为空数组[] */
				{					
					id: {INT},		/* 获得item类型 */
					num: {INT}		/* 获得item数量 */
				},
				...
			],
			nextPvpBands: [					/* pvp匹配到的3个队伍，如果没有为[] */
				{							/* pvp队伍信息 */
					userId:{INT}			/* 对手用户id */
					userName:{STRING}		/* 对手用户名 */
					userLevel:{INT}			/* 对手的等级 */
					formation:{INT}			/* 阵形 */
					cards:[					/* 对手band */
						{					/* band中的成员，如果位置为空，则为null */
						    protoId:{INT},			/* 原型id（卡片类型） */
						    level:{INT},			/* 等级 */
						    hp:{INT},				/* hp总和 */
						    atk:{INT},				/* atk总和 */
						    def:{INT},				/* def总和 */
						    wis:{INT},				/* wis总和 */
						    agi:{INT},				/* agi总和 */
						    skill1Id:{INT},			/* 技能1 id */
						    skill1Level:{INT},		/* 技能1 等级 */
						    skill2Id:{INT},			/* 技能2 id */
						    skill2Level:{INT},		/* 技能2 等级 */
						    skill3Id:{INT},			/* 技能3 id */
						    skill3Level:{INT}		/* 技能3 等级 */
						},
						...
					]
				},
				...
			]
		}




