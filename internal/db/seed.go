package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"

	"github.com/robertgouveia/social/internal/store"
)

var usernames = []string{
	"CoolCat", "SkyWalker", "NeonNinja", "FastFalcon", "PixelPanda", "ShadowSeeker", "ElectricEagle", "AcePilot", "StormSurfer", "BlazeTrail",
	"MysticMuse", "CyberWolf", "ThunderBolt", "AquaArcher", "LunarLion", "IronClad", "QuantumQuest", "EchoNight", "SolarScribe", "DreamDrifter",
	"FireFury", "GravityGuru", "HexHawk", "OrbitOracle", "SwiftSparrow", "InfernoIvy", "FrostFox", "ZenZero", "TidalTiger", "CloudCrafter",
	"WindWanderer", "SteelShade", "CrimsonCrow", "SparkPhoenix", "NovaNoble", "StoneSculptor", "GlitterGale", "HorizonHawk", "AstralArrow", "EchoElf",
	"WildWhisper", "DynamoDuck", "KnightVigil", "FableForge", "StormStriker", "RuneRaven", "SwiftSage", "PulsePilot", "ChillCrusader", "EpicEagle",
}

var postTitles = []string{
	"Exploring the Unknown", "Top 10 Travel Hacks", "Mastering the Art of Focus", "Healthy Habits for Life",
	"Secrets of Great Leaders", "Boost Your Productivity", "DIY Home Makeover Tips", "Learning to Code 101",
	"Mindfulness for Beginners", "How to Stay Motivated", "Creating a Morning Routine", "Budget Travel Guide",
	"Top Books to Read", "Quick & Easy Recipes", "Guide to Personal Finance", "Fitness at Home",
	"Photography Tips & Tricks", "Mastering Time Management", "Gardening for Beginners", "Overcoming Procrastination",
}

var postContents = []string{
	"Discover the thrill of venturing into the unknown and the rewards it brings.",
	"Planning a trip? Here are the top 10 travel hacks to make your journey smoother.",
	"Learn key techniques to improve your focus and achieve your goals efficiently.",
	"Building healthy habits can be life-changing—here's a guide to get started.",
	"Uncover what makes great leaders tick, and learn how you can apply it to your life.",
	"Boost your productivity with these proven tips and tools for a more effective day.",
	"Give your home a fresh look with these DIY makeover tips—easy and affordable!",
	"Want to start coding? Here’s a beginner’s guide to kickstart your journey in tech.",
	"Curious about mindfulness? Learn the basics and see how it can transform your days.",
	"Need motivation? Here are simple strategies to keep your energy and enthusiasm high.",
	"Create a morning routine that sets you up for success with these easy steps.",
	"Traveling on a budget doesn’t have to be hard—here’s how to explore affordably.",
	"Looking for a good read? Check out this list of must-read books across genres.",
	"Cooking at home is easier than ever with these quick and easy recipes for busy days.",
	"Take control of your finances with practical advice on budgeting and saving.",
	"Get fit without leaving home with these simple exercises and routines.",
	"Elevate your photography skills with these essential tips for beginners and pros.",
	"Master time management with these techniques to make the most of each day.",
	"Start your own garden with this beginner's guide, from plants to tools.",
	"Break free from procrastination with strategies to help you take action today.",
}

var tags = []string{
	"Adventure", "Travel", "Productivity", "Health", "Leadership",
	"DIY", "Coding", "Mindfulness", "Motivation", "Routine",
	"Budgeting", "Books", "Recipes", "Finance", "Fitness",
	"Photography", "Time Management", "Gardening", "Self-Improvement", "Lifestyle",
}

var comments = []string{
	"Great tips! I always wanted to explore more.",
	"These travel hacks are a lifesaver. Thanks for sharing!",
	"I struggle with focus; I'll definitely try these techniques.",
	"Healthy habits make such a difference in my life. Love this!",
	"Interesting perspective on leadership; very inspiring!",
	"Productivity tips are always welcome! Can't wait to try them.",
	"I love DIY projects! Can’t wait to start my makeover.",
	"Coding is intimidating, but this guide makes it feel achievable!",
	"Mindfulness has really changed my approach to stress. Thank you!",
	"I could use some motivation today; this was perfect!",
	"Morning routines are crucial! I've been wanting to improve mine.",
	"Budget travel ideas are gold! I need this for my next trip.",
	"Excited to check out these book recommendations!",
	"Quick recipes are a lifesaver for busy weeknights. Thank you!",
	"Great advice on finance; I’m working on my budget.",
	"I never thought fitness could be this fun at home!",
	"Photography tips are always appreciated; can't wait to learn!",
	"Time management is key! I need to implement these tips.",
	"Gardening seems fun! I'm thinking of starting small.",
	"Procrastination is my biggest enemy. This is motivating!",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding Complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			//modulo allows for looping (5 % 5 = 0 -- first title)
			Username: usernames[i%len(usernames)] + strconv.Itoa(rand.IntN(1000)),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.IntN(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   postTitles[rand.IntN(len(postTitles))],
			Content: postContents[rand.IntN(len(postContents))],
			Tags:    generateTags(3),
		}
	}

	return posts
}

func generateTags(maxNum int) []string {
	randNum := rand.IntN(maxNum)
	res := make([]string, randNum)
	for i := 0; i < randNum; i++ {
		res = append(res, tags[rand.IntN(len(tags))])
	}

	return res
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}

	return cms
}
